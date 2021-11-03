package v1

import (
	"bytes"
	json2 "encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/docker"
	upnp2 "github.com/IceWhaleTech/CasaOS/pkg/upnp"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	ip_helper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/ip_helper"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	port2 "github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/random"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/IceWhaleTech/CasaOS/service/docker_base"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Duration(time.Second * 5),
}

//打开docker的terminal
func DockerTerminal(c *gin.Context) {
	col := c.DefaultQuery("cols", "100")
	row := c.DefaultQuery("rows", "30")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	defer conn.Close()
	container := c.Param("id")
	hr, err := service.Exec(container, row, col)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()
	go func() {
		docker.WsWriterCopy(hr.Conn, conn)
	}()
	docker.WsReaderCopy(conn, hr.Conn)
}

//打开本机的ssh接口
func WsSsh(c *gin.Context) {
	wsConn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

	var logBuff = new(bytes.Buffer)
	quitChan := make(chan bool, 3)
	user := ""
	password := ""
	var login int = 1
	cols, _ := strconv.Atoi(c.DefaultQuery("cols", "200"))
	rows, _ := strconv.Atoi(c.DefaultQuery("rows", "32"))
	var client *ssh.Client
	for login != 0 {

		var err error

		wsConn.WriteMessage(websocket.TextMessage, []byte("login:"))
		user = docker.ReceiveWsMsgUser(wsConn, logBuff)
		wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		wsConn.WriteMessage(websocket.TextMessage, []byte("password:"))
		password = docker.ReceiveWsMsgPassword(wsConn, logBuff)
		wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		client, err = docker.NewSshClient(user, password)

		if err != nil && client == nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		} else {
			login = 0
		}

	}
	if client != nil {
		defer client.Close()
	}

	ssConn, _ := docker.NewSshConn(cols, rows, client)
	defer ssConn.Close()

	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan

}

//安装进度推送
func SpeedPush(c *gin.Context) {
	//token := c.Query("token")
	//if len(token) == 0 || token != config.UserInfo.Token {
	//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR_AUTH_TOKEN, Message: oasis_err2.GetMsg(oasis_err2.ERROR_AUTH_TOKEN)})
	//	return
	//}

	//ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	//defer ws.Close()
	//
	//for {
	//	select {
	//	case msg := <-WSMSG:
	//		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintln(msg)))
	//	}
	//}
}

// @Summary 安装app(该接口需要post json数据)
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path int true "id"
// @Param  port formData int true "主端口"
// @Param  tcp formData string false "tcp端口"
// @Param  udp formData string false "udp端口"
// @Param  env formData string false "环境变量"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/install/{id} [post]
func InstallApp(c *gin.Context) {
	appId := c.Param("id")
	var appInfo model.ServerAppList
	m := model.CustomizationPostData{}
	c.BindJSON(&m)

	const CUSTOM = "custom"
	var dockerImage string
	var dockerImageVersion string
	//检查端口
	if len(m.PortMap) > 0 && m.PortMap != "0" {
		//c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		portMap, _ := strconv.Atoi(m.PortMap)
		if !port2.IsPortAvailable(portMap, "tcp") {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + m.PortMap})
			return
		}
	}
	//if len(m.Port) == 0 || m.Port == "0" {
	//	m.Port = m.PortMap
	//}

	imageArr := strings.Split(m.Image, ":")
	if len(imageArr) == 2 {
		dockerImage = imageArr[0]
		dockerImageVersion = imageArr[1]
	} else {
		dockerImage = m.Image
		dockerImageVersion = "latest"
	}
	if m.Origin != "custom" {
		appInfo = service.MyService.OAPI().GetServerAppInfo(appId)

	} else {

		appInfo.Title = m.Label
		appInfo.Description = m.Description
		appInfo.Icon = m.Icon
		appInfo.ScreenshotLink = model.Strings{}
		appInfo.NetworkModel = m.NetworkModel
		appInfo.Tags = model.Strings{}
		appInfo.Tagline = ""
		appInfo.Index = m.Index

	}

	for _, u := range m.Ports {

		if u.Protocol == "udp" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "tcp" {

			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "both" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		}

	}

	for _, device := range m.Devices {
		if file.CheckNotExist(device.Path) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.DEVICE_NOT_EXIST, Message: device.Path + "," + oasis_err2.GetMsg(oasis_err2.DEVICE_NOT_EXIST)})
			return
		}

	}
	//restart := c.PostForm("restart") //always 总是重启,   unless-stopped 除非用户手动停止容器，否则总是重新启动,    on-failure:仅当容器退出代码非零时重新启动
	//if len(restart) > 0 {
	//
	//}
	//
	//privileged := c.PostForm("privileged") //是否处于特权模式
	//if len(privileged) > 0 {
	//
	//}
	id := uuid.NewV4().String()

	var relyMap = make(map[string]string)
	go func() {
		installLog := model2.AppNotify{}
		installLog.State = 0
		installLog.CustomId = id
		installLog.Message = "installing rely"
		installLog.Type = types.NOTIFY_TYPE_UNIMPORTANT
		installLog.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		service.MyService.Notify().AddLog(installLog)
		if m.Origin != "custom" {
			for _, plugin := range appInfo.Plugins {
				if plugin == "mysql" {
					mid := uuid.NewV4().String()
					mc := docker_base.MysqlConfig{}
					mc.DataBasePassword = random.RandomString(6, false)
					mc.DataBaseDB = appInfo.Title
					mc.DataBaseUser = "root"
					mc.DataBasePort = "3306"
					mysqlContainerId, err := docker_base.MysqlCreate(mc, mid, m.CpuShares, m.Memory)
					if len(mysqlContainerId) > 0 && err == nil {

						mc.DataBaseHost = mid

						m.Envs = docker_base.MysqlFilter(mc, m.Envs)

						rely := model2.RelyDBModel{}
						rely.Type = types.RELY_TYPE_MYSQL
						rely.ContainerId = mysqlContainerId
						rely.CustomId = mid
						rely.ContainerCustomId = id
						var msqlConfig model2.MysqlConfigs

						//结构体转换
						copier.Copy(&msqlConfig, &mc)
						rely.Config = msqlConfig
						service.MyService.Rely().Create(rely)

						relyMap["mysql"] = mid

					} else {
						docker_base.MysqlDelete(mysqlContainerId)
						installLog.State = 0
						installLog.Message = err.Error()
						service.MyService.Notify().UpdateLog(installLog)
					}
				}
			}
		}

		installLog.Message = "pulling"
		service.MyService.Notify().UpdateLog(installLog)

		// step：下载镜像
		err := service.MyService.Docker().DockerPullImage(dockerImage+":"+dockerImageVersion, installLog)
		if err != nil {
			installLog.State = 0
			installLog.Message = err.Error()
			installLog.Type = types.NOTIFY_TYPE_ERROR
			service.MyService.Notify().UpdateLog(installLog)
			return
		}

		for !service.MyService.Docker().IsExistImage(dockerImage + ":" + dockerImageVersion) {
			time.Sleep(time.Second)
		}

		//if  {

		//}

		//step：创建容器
		// networkName, err := service.MyService.Docker().GetNetWorkNameByNetWorkID(appInfo.NetworkModel)
		// if err != nil {
		// 	//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":80}", 100)
		// 	installLog.State = 0
		// 	installLog.Speed = 75
		// 	installLog.Type = types.NOTIFY_TYPE_ERROR
		// 	installLog.Message = err.Error()
		// 	service.MyService.Notify().UpdateLog(installLog)
		// 	return
		// }
		containerId, err := service.MyService.Docker().DockerContainerCreate(dockerImage+":"+dockerImageVersion, id, m, appInfo.NetworkModel)
		installLog.Name = appInfo.Title
		installLog.Icon = appInfo.Icon
		if err != nil {
			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":80}", 100)
			installLog.State = 0
			installLog.Type = types.NOTIFY_TYPE_ERROR
			installLog.Message = err.Error()
			service.MyService.Notify().UpdateLog(installLog)
			return
		} else {
			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"starting\",\"speed\":80}", 100)
			installLog.Message = "starting"
			service.MyService.Notify().UpdateLog(installLog)
		}

		//		echo -e "hellow\nworld" >>

		//step：启动容器
		err = service.MyService.Docker().DockerContainerStart(id)
		if err != nil {
			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":90}", 100)
			installLog.State = 0
			installLog.Type = types.NOTIFY_TYPE_ERROR
			installLog.Message = err.Error()
			service.MyService.Notify().UpdateLog(installLog)
			return
		} else {
			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"setting upnp\",\"speed\":90}", 100)
			if m.Origin != CUSTOM {
				installLog.Message = "setting upnp"
			} else {
				installLog.Message = "nearing completion"
			}
			service.MyService.Notify().UpdateLog(installLog)
		}

		if m.Origin != CUSTOM {
			//step:启动upnp
			if m.EnableUPNP {
				upnp, err := upnp2.Gateway()
				if err == nil {
					for _, p := range m.Ports {
						if p.Protocol == "udp" {
							upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
							upnp.LocalHost = ip_helper2.GetLoclIp()
							tComment, _ := strconv.Atoi(p.CommendPort)
							upnp.AddPortMapping(tComment, tComment, "UDP")
							time.Sleep(time.Millisecond * 200)
						} else if p.Protocol == "tcp" {
							upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
							upnp.LocalHost = ip_helper2.GetLoclIp()
							tComment, _ := strconv.Atoi(p.CommendPort)
							upnp.AddPortMapping(tComment, tComment, "TCP")
							time.Sleep(time.Millisecond * 200)
						} else if p.Protocol == "both" {

							upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
							upnp.LocalHost = ip_helper2.GetLoclIp()
							tComment, _ := strconv.Atoi(p.CommendPort)
							upnp.AddPortMapping(tComment, tComment, "UDP")
							time.Sleep(time.Millisecond * 200)

							upnp.AddPortMapping(tComment, tComment, "TCP")
							time.Sleep(time.Millisecond * 200)
						}

					}

				}
				if err != nil {
					//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":95}", 100)
					installLog.State = 0
					installLog.Type = types.NOTIFY_TYPE_ERROR
					installLog.Message = err.Error()
					service.MyService.Notify().UpdateLog(installLog)
				} else {
					//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"checking\",\"speed\":95}", 100)
					installLog.Message = "checking"
					service.MyService.Notify().UpdateLog(installLog)
				}
			}
		}

		//step: 启动成功     检查容器状态确认启动成功
		container, err := service.MyService.Docker().DockerContainerInfo(id)
		if err != nil && container.ContainerJSONBase.State.Running {
			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":100}", 100)
			installLog.State = 0
			installLog.Type = types.NOTIFY_TYPE_ERROR
			installLog.Message = err.Error()
			service.MyService.Notify().UpdateLog(installLog)
			return
		} else {
			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"installed\",\"speed\":100}", 100)
			installLog.Message = "installed"
			service.MyService.Notify().UpdateLog(installLog)
		}

		rely := model.MapStrings{}

		copier.Copy(&rely, &relyMap)
		portsStr, _ := json2.Marshal(m.Ports)
		envsStr, _ := json2.Marshal(m.Envs)
		volumesStr, _ := json2.Marshal(m.Volumes)
		devicesStr, _ := json2.Marshal(m.Devices)
		//step: 保存数据到数据库
		md := model2.AppListDBModel{
			CustomId: id,
			Title:    appInfo.Title,
			//ScreenshotLink: appInfo.ScreenshotLink,
			Slogan:      appInfo.Tagline,
			Description: appInfo.Description,
			//Tags:           appInfo.Tags,
			Icon:        appInfo.Icon,
			Version:     dockerImageVersion,
			ContainerId: containerId,
			Image:       dockerImage,
			Index:       appInfo.Index,
			//Port:           m.Port,
			PortMap:    m.PortMap,
			Label:      m.Label,
			EnableUPNP: m.EnableUPNP,
			Ports:      string(portsStr),
			Envs:       string(envsStr),
			Volumes:    string(volumesStr),
			Position:   m.Position,
			NetModel:   appInfo.NetworkModel,
			Restart:    m.Restart,
			CpuShares:  m.CpuShares,
			Memory:     m.Memory,
			Devices:    string(devicesStr),
			//Rely:       rely,
			Origin:    m.Origin,
			CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
			UpdatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		}
		//if appInfo.NetworkModel == "host" {
		//	m.PortMap = m.Port
		//}
		service.MyService.App().SaveContainer(md)

	}()

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: id})

}

//// @Summary 自定义安装app(该接口需要post json数据)
//// @Produce  application/json
//// @Accept application/json
//// @Tags app
//// @Param  id path int true "id"
//// @Param  port formData int true "主端口"
//// @Param  tcp formData string false "tcp端口"
//// @Param  udp formData string false "udp端口"
//// @Param  env formData string false "环境变量"
//// @Security ApiKeyAuth
//// @Success 200 {string} string "ok"
//// @Router /app/install/{id} [post]
//func CustomInstallApp(c *gin.Context) {
//	//appId := c.Param("id")
//	//	appInfo := service.MyService.App().GetServerAppInfo(appId)
//
//	m := model.CustomizationPostData{}
//	c.BindJSON(&m)
//	//检查端口
//	if len(m.PortMap) == 0 || m.PortMap == "0" {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	if len(m.Port) == 0 || m.Port == "0" {
//		m.Port = m.PortMap
//	}
//
//	portMap, _ := strconv.Atoi(m.PortMap)
//	if !port2.IsPortAvailable(portMap, "tcp") {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + m.PortMap})
//		return
//	}
//
//	for _, u := range m.Udp {
//		t, _ := strconv.Atoi(u.CommendPort)
//		if !port2.IsPortAvailable(t, "udp") {
//			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
//			return
//		}
//	}
//
//	for _, t := range m.Tcp {
//		te, _ := strconv.Atoi(t.CommendPort)
//		if !port2.IsPortAvailable(te, "tcp") {
//			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + t.CommendPort})
//			return
//		}
//	}
//
//	//restart := c.PostForm("restart") //always 总是重启,   unless-stopped 除非用户手动停止容器，否则总是重新启动,    on-failure:仅当容器退出代码非零时重新启动
//	//if len(restart) > 0 {
//	//
//	//}
//	//
//	//privileged := c.PostForm("privileged") //是否处于特权模式
//	//if len(privileged) > 0 {
//	//
//	//}
//
//	err := service.MyService.Docker().DockerPullImage(m.Image)
//	if err != nil {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PULL_IMAGE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PULL_IMAGE_ERROR)})
//	}
//
//	id := uuid.NewV4().String()
//
//	var relyMap = make(map[string]string)
//	go func() {
//		installLog := model2.AppNotify{}
//		installLog.CustomId = id
//		installLog.State = 0
//		installLog.Message = "installing rely"
//		installLog.Speed = 30
//		installLog.CreatedAt = time.Now()
//		installLog.UpdatedAt = time.Now()
//		service.MyService.Notify().AddLog(installLog)
//
//		for !service.MyService.Docker().IsExistImage(m.Image) {
//			time.Sleep(time.Second)
//		}
//
//		installLog.Speed = 50
//		installLog.Message = "pulling"
//		service.MyService.Notify().UpdateLog(installLog)
//		// step：下载镜像
//
//		var cpd model.PostData
//		copier.Copy(&cpd, &m)
//		//step：创建容器
//		containerId, err := service.MyService.Docker().DockerContainerCreate(m.Image, id, cpd, m.NetworkModel, m.Image, "custom")
//		installLog.ContainerId = containerId
//		if err != nil {
//			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":80}", 100)
//			installLog.State = 0
//			installLog.Speed = 80
//			installLog.Message = err.Error()
//			service.MyService.Notify().UpdateLog(installLog)
//			return
//		} else {
//			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"starting\",\"speed\":80}", 100)
//			installLog.Speed = 80
//			installLog.Message = "starting"
//			service.MyService.Notify().UpdateLog(installLog)
//		}
//
//		//step：启动容器
//		err = service.MyService.Docker().DockerContainerStart(id)
//		if err != nil {
//			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":90}", 100)
//			installLog.State = 0
//			installLog.Speed = 90
//			installLog.Message = err.Error()
//			service.MyService.Notify().UpdateLog(installLog)
//			return
//		} else {
//			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"setting upnp\",\"speed\":90}", 100)
//			installLog.Speed = 90
//			installLog.Message = "setting upnp"
//			service.MyService.Notify().UpdateLog(installLog)
//		}
//
//		//step: 启动成功     检查容器状态确认启动成功
//		containerStatus, err := service.MyService.Docker().DockerContainerInfo(id)
//		if err != nil && containerStatus.ContainerJSONBase.State.Running {
//			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":100}", 100)
//			installLog.State = 0
//			installLog.Speed = 100
//			installLog.Message = err.Error()
//			service.MyService.Notify().UpdateLog(installLog)
//			return
//		} else {
//			//service.MyService.Redis().Set(id, "{\"id\":\""+id+"\",\"state\":true,\"message\":\"installed\",\"speed\":100}", 100)
//			installLog.Speed = 100
//			installLog.Message = "installed"
//			service.MyService.Notify().UpdateLog(installLog)
//		}
//
//		rely := model.MapStrings{}
//
//		copier.Copy(&rely, &relyMap)
//
//		//step: 保存数据到数据库
//		md := model2.AppListDBModel{
//			CustomId: id,
//			Title:    m.Label,
//			//			ScreenshotLink: []string,
//			Slogan:      "",
//			Description: m.Description,
//			//			Tags:           ,
//			Icon:        m.Icon,
//			Version:     m.Image,
//			ContainerId: containerId,
//			Image:       m.Image,
//			Index:       "",
//			Port:        m.Port,
//			PortMap:     m.PortMap,
//			Label:       m.Label,
//			EnableUPNP:  m.EnableUPNP,
//			UdpPorts:    m.Udp,
//			TcpPorts:    m.Tcp,
//			Envs:        m.Envs,
//			Volumes:     m.Volumes,
//			Position:    m.Position,
//			NetModel:    m.NetworkModel,
//			Restart:     m.Restart,
//			CpuShares:   m.CpuShares,
//			Memory:      m.Memory,
//			Devices:     m.Devices,
//			Rely:        rely,
//			Origin:      "custom",
//		}
//		if m.NetworkModel == "host" {
//			m.PortMap = m.Port
//		}
//		service.MyService.App().SaveContainer(md)
//
//	}()
//
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: id})
//
//}

// @Summary 卸载app
// @Produce  application/json
// @Accept multipart/form-data
// @Tags app
// @Param  id path string true "容器id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/uninstall/{id} [delete]
func UnInstallApp(c *gin.Context) {
	appId := c.Param("id")

	if len(appId) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.App().GetUninstallInfo(appId)

	//step：停止容器
	err := service.MyService.Docker().DockerContainerStop(appId)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.UNINSTALL_APP_ERROR, Message: oasis_err2.GetMsg(oasis_err2.UNINSTALL_APP_ERROR), Data: err.Error()})
		return
	}

	//step：删除容器
	err = service.MyService.Docker().DockerContainerRemove(appId)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.UNINSTALL_APP_ERROR, Message: oasis_err2.GetMsg(oasis_err2.UNINSTALL_APP_ERROR), Data: err.Error()})
		return
	}

	//存在镜像正在使用的情况
	// step：删除镜像
	service.MyService.Docker().DockerImageRemove(info.Image + ":" + info.Version)

	//step: 删除本地数据
	service.MyService.App().RemoveContainerById(appId)
	if info.Origin != "custom" {

		//step: 删除文件夹
		vol := gjson.Get(info.Volumes, "#.host")
		for _, v := range vol.Array() {

			service.MyService.App().DelAppConfigDir(appId, v.String())
		}

		//step: 删除install log
		service.MyService.Notify().DelLog(appId)

		//	for k, v := range info.Rely {
		//
		//		if k == "mysql" {
		//			docker_base.MysqlDelete(v)
		//			service.MyService.Rely().Delete(v)
		//		}
		//	}

		//if info.EnableUPNP {
		//	upnp, err := upnp2.Gateway()
		//	if err == nil {
		//		for _, p := range info.Ports {
		//			if p.Protocol == "udp" {
		//				upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//				upnp.LocalHost = ip_helper2.GetLoclIp()
		//				tComment, _ := strconv.Atoi(p.CommendPort)
		//				upnp.DelPortMapping(tComment, "UDP")
		//				time.Sleep(time.Millisecond * 200)
		//			} else if p.Protocol == "tcp" {
		//				upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//				upnp.LocalHost = ip_helper2.GetLoclIp()
		//				tComment, _ := strconv.Atoi(p.CommendPort)
		//				upnp.DelPortMapping(tComment, "TCP")
		//				time.Sleep(time.Millisecond * 200)
		//			} else if p.Protocol == "both" {
		//				upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//				upnp.LocalHost = ip_helper2.GetLoclIp()
		//				tComment, _ := strconv.Atoi(p.CommendPort)
		//				upnp.DelPortMapping(tComment, "UDP")
		//
		//				upnp.DelPortMapping(tComment, "TCP")
		//				time.Sleep(time.Millisecond * 200)
		//			}
		//		}
		//	}
		//}
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})

}

// @Summary 修改app状态
// @Produce  application/json
// @Accept multipart/form-data
// @Tags app
// @Param  id path string true "appid"
// @Param  state query string false "是否停止 start stop restart"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/state/{id} [put]
func ChangAppState(c *gin.Context) {
	appId := c.Param("id")
	state := c.DefaultPostForm("state", "stop")
	var err error
	if state == "stop" {
		err = service.MyService.Docker().DockerContainerStop(appId)
	} else if state == "start" {
		err = service.MyService.Docker().DockerContainerStart(appId)
	} else if state == "restart" {
		err = service.MyService.Docker().DockerContainerStop(appId)
		err = service.MyService.Docker().DockerContainerStart(appId)
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	info, err := service.MyService.App().GetContainerInfo(appId)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info.State})
}

// @Summary 查看容器日志
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "appid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/logs/{id} [get]
func ContainerLog(c *gin.Context) {
	appId := c.Param("id")
	log, _ := service.MyService.Docker().DockerContainerLog(appId)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: log})
}

// @Summary 获取安装进度
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "容器id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/speed/{id} [get]
func GetInstallSpeed(c *gin.Context) {
	id := c.Param("id")
	b := service.MyService.Notify().GetLog(id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: b})
}

// @Summary 获取容器状态
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "容器id"
// @Param  type query string false "type=1"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/state/{id} [get]
func GetContainerState(c *gin.Context) {
	id := c.Param("id")
	t := c.DefaultQuery("type", "0")
	containerInfo, e := service.MyService.App().GetSimpleContainerInfo(id)
	if e != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: e.Error()})
		return
	}

	var data = make(map[string]interface{})

	data["state"] = containerInfo.State

	if t == "1" {
		appInfo := service.MyService.App().GetAppDBInfo(id)
		data["app"] = appInfo
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: data})
}

// @Summary 更新设置
// @Produce  application/json
// @Accept multipart/form-data
// @Tags app
// @Param  id path string true "容器id"
// @Param  shares formData string false "cpu权重"
// @Param  mem formData string false "内存大小MB"
// @Param  restart formData string false "重启策略"
// @Param  label formData string false "应用名称"
// @Param  position formData bool true "是否放到首页"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/update/{id}/setting [put]
func UpdateSetting(c *gin.Context) {
	id := c.Param("id")
	const CUSTOM = "custom"
	m := model.CustomizationPostData{}
	c.BindJSON(&m)

	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	var cpd model.CustomizationPostData

	copier.Copy(&cpd, &m)

	appInfo := service.MyService.App().GetAppDBInfo(id)

	var containerId string
	containerId = appInfo.ContainerId

	service.MyService.Docker().DockerContainerStop(id)

	portMap, _ := strconv.Atoi(m.PortMap)
	if !port2.IsPortAvailable(portMap, "tcp") {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + m.PortMap})
		return
	}

	for _, u := range m.Ports {

		if u.Protocol == "udp" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "tcp" {
			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "both" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}

			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		}

	}

	//如果容器端口均未修改,这不进行处理
	portsStr, _ := json2.Marshal(m.Ports)

	envsStr, _ := json2.Marshal(m.Envs)
	volumesStr, _ := json2.Marshal(m.Volumes)
	devicesStr, _ := json2.Marshal(m.Devices)
	if !reflect.DeepEqual(string(portsStr), appInfo.Ports) || !reflect.DeepEqual(string(envsStr), appInfo.Envs) || !reflect.DeepEqual(string(volumesStr), appInfo.Volumes) || m.PortMap != appInfo.PortMap || m.NetworkModel != appInfo.NetModel {

		var newUUid = uuid.NewV4().String()
		var err error

		// networkName, err := service.MyService.Docker().GetNetWorkNameByNetWorkID(appInfo.NetModel)
		// if err != nil {
		// 	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
		// 	return
		// }
		containerId, err = service.MyService.Docker().DockerContainerCreate(appInfo.Image+":"+appInfo.Version, newUUid, cpd, m.NetworkModel)

		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
			return
		}

		err = service.MyService.Docker().DockerContainerRemove(id)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
			return
		}

		service.MyService.Docker().DockerContainerUpdateName(appInfo.CustomId, newUUid)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
			return
		}

	} else if !reflect.DeepEqual(string(devicesStr), appInfo.Devices) || m.CpuShares != appInfo.CpuShares || m.Memory != appInfo.Memory || m.Restart != appInfo.Restart {
		service.MyService.Docker().DockerContainerUpdate(cpd, id)
	}

	err := service.MyService.Docker().DockerContainerStart(id)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
		return
	}
	//更新upnp
	if m.Origin != CUSTOM {
		//if appInfo.EnableUPNP != appInfo.EnableUPNP {
		//	if appInfo.EnableUPNP {
		//		upnp, err := upnp2.Gateway()
		//		if err == nil {
		//
		//			for _, p := range appInfo.Ports {
		//				if p.Protocol == "udp" {
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.AddPortMapping(tComment, tComment, "UDP")
		//					time.Sleep(time.Millisecond * 200)
		//				} else if p.Protocol == "tcp" {
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.AddPortMapping(tComment, tComment, "TCP")
		//					time.Sleep(time.Millisecond * 200)
		//				} else if p.Protocol == "both" {
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.AddPortMapping(tComment, tComment, "UDP")
		//					time.Sleep(time.Millisecond * 200)
		//
		//					upnp.AddPortMapping(tComment, tComment, "TCP")
		//					time.Sleep(time.Millisecond * 200)
		//				}
		//			}
		//		}
		//	} else {
		//		upnp, err := upnp2.Gateway()
		//		if err == nil {
		//			for _, p := range appInfo.Ports {
		//				if p.Protocol == "udp" {
		//
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.DelPortMapping(tComment, "UDP")
		//					time.Sleep(time.Millisecond * 200)
		//				} else if p.Protocol == "tcp" {
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.DelPortMapping(tComment, "TCP")
		//					time.Sleep(time.Millisecond * 200)
		//				} else if p.Protocol == "both" {
		//					upnp.CtrlUrl = upnp2.GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		//					upnp.LocalHost = ip_helper2.GetLoclIp()
		//					tComment, _ := strconv.Atoi(p.CommendPort)
		//					upnp.DelPortMapping(tComment, "UDP")
		//					time.Sleep(time.Millisecond * 200)
		//
		//					upnp.DelPortMapping(tComment, "TCP")
		//					time.Sleep(time.Millisecond * 200)
		//				}
		//			}
		//		}
		//	}
		//}
	}

	appInfo.ContainerId = containerId
	appInfo.PortMap = m.PortMap
	appInfo.Label = m.Label
	appInfo.Index = m.Index
	appInfo.Ports = string(portsStr)
	appInfo.Envs = string(envsStr)
	appInfo.Icon = m.Icon
	appInfo.Volumes = string(volumesStr)
	appInfo.Devices = string(devicesStr)
	appInfo.NetModel = m.NetworkModel
	appInfo.Position = m.Position
	appInfo.EnableUPNP = m.EnableUPNP
	appInfo.Restart = m.Restart
	appInfo.Memory = m.Memory
	appInfo.CpuShares = m.CpuShares
	appInfo.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
	service.MyService.App().UpdateApp(appInfo)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary 获取容器详情
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "appid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/info/{id} [get]
func ContainerInfo(c *gin.Context) {
	appId := c.Param("id")
	appInfo := service.MyService.App().GetAppDBInfo(appId)
	containerInfo, _ := service.MyService.Docker().DockerContainerStats(appId)
	var cpuModel = "arm"
	if cpu := service.MyService.ZiMa().GetCpuInfo(); len(cpu) > 0 {
		if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "intel") > 0 {
			cpuModel = "intel"
		} else if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "amd") > 0 {
			cpuModel = "amd"
		}
	}

	info, err := service.MyService.Docker().DockerContainerInfo(appId)
	if err != nil {
		//todo 需要自定义错误
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: err.Error()})
		return
	}
	con := struct {
		Status    string `json:"status"`
		StartedAt string `json:"started_at"`
		CPUShares int64  `json:"cpu_shares"`
		Memory    int64  `json:"memory"`
		Restart   string `json:"restart"`
	}{Status: info.State.Status, StartedAt: info.State.StartedAt, CPUShares: info.HostConfig.CPUShares, Memory: info.HostConfig.Memory >> 20, Restart: info.HostConfig.RestartPolicy.Name}
	data := make(map[string]interface{}, 5)
	data["app"] = appInfo
	data["cpu"] = cpuModel
	data["memory"] = service.MyService.ZiMa().GetMemInfo().Total
	data["container"] = json2.RawMessage(containerInfo)
	data["info"] = con
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: data})
}

// @Summary 获取安装所需要的数据
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/install/config [get]
func GetDockerInstallConfig(c *gin.Context) {
	networks := service.MyService.Docker().DockerNetworkModelList()
	data := make(map[string]interface{}, 2)
	list := []map[string]string{}
	for _, network := range networks {
		if network.Driver != "null" {
			list = append(list, map[string]string{"name": network.Name, "driver": network.Driver, "id": network.ID})
		}
	}
	data["networks"] = list
	data["memory"] = service.MyService.ZiMa().GetMemInfo()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: data})
}

// @Summary 获取依赖数据
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "rely id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/rely/{id}/info [get]
func ContainerRelyInfo(c *gin.Context) {
	id := c.Param("id")
	appInfo := service.MyService.Rely().GetInfo(id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: appInfo})
}

// @Summary 获取可更新数据
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "appid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/update/{id}/info [get]
func ContainerUpdateInfo(c *gin.Context) {
	appId := c.Param("id")
	appInfo := service.MyService.App().GetAppDBInfo(appId)

	info, err := service.MyService.Docker().DockerContainerInfo(appId)
	if err != nil {
		//todo 需要自定义错误
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: err.Error()})
		return
	}
	var port model.PortArray
	json2.Unmarshal([]byte(appInfo.Ports), &port)

	var envs model.EnvArray
	json2.Unmarshal([]byte(appInfo.Envs), &envs)

	var vol model.PathArray
	json2.Unmarshal([]byte(appInfo.Volumes), &vol)

	var dir model.PathArray
	json2.Unmarshal([]byte(appInfo.Devices), &dir)

	//volumesStr, _ := json2.Marshal(m.Volumes)
	//devicesStr, _ := json2.Marshal(m.Devices)
	m := model.CustomizationPostData{}
	m.Index = appInfo.Index
	m.Icon = appInfo.Icon
	m.Ports = port
	m.Image = appInfo.Image + ":" + appInfo.Version
	m.Origin = appInfo.Origin
	m.NetworkModel = appInfo.NetModel
	m.Description = appInfo.Description
	m.Label = appInfo.Label
	m.PortMap = appInfo.PortMap
	m.Devices = dir //appInfo.Devices
	m.Envs = envs
	m.Memory = info.HostConfig.Memory >> 20
	m.CpuShares = info.HostConfig.CPUShares
	m.Volumes = vol //appInfo.Volumes
	m.Restart = info.HostConfig.RestartPolicy.Name
	m.EnableUPNP = appInfo.EnableUPNP
	m.Position = appInfo.Position

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: m})
}

////准备安装(暂时不需要)
//func ReadyInstall(c *gin.Context) {
//	_, tcp, udp := service.MyService.GetManifestJsonByRepo()
//	data := make(map[string]interface{}, 2)
//	if t := gjson.Parse(tcp).Array(); len(t) > 0 {
//		//tcpList := []model.TcpPorts{}
//		//e := json2.Unmarshal([]byte(tcp), tcpList)
//		//if e!=nil {
//		//	return
//		//}
//		//for _, port := range tcpList {
//		//	if port.ContainerPort>0&&port.ExtranetPort {
//		//
//		//	}
//		//}
//		var inarr []interface{}
//		for _, result := range t {
//
//			var p int
//			ok := true
//			for ok {
//				p, _ = port.GetAvailablePort()
//				ok = !port.IsPortAvailable(p)
//			}
//			pm := model.PortMap{gjson.Get(result.Raw, "container_port").Int(), p}
//			inarr = append(inarr, pm)
//		}
//		data["tcp"] = inarr
//	}
//	if u := gjson.Parse(udp).Array(); len(u) > 0 {
//		//udpList := []model.UdpPorts{}
//		//e := json2.Unmarshal([]byte(udp), udpList)
//		//if e != nil {
//		//	return
//		//}
//		var inarr []model.PortMap
//		for _, result := range u {
//			var p int
//			ok := true
//			for ok {
//				p, _ = port.GetAvailablePort()
//				ok = !port.IsPortAvailable(p)
//			}
//			pm := model.PortMap{gjson.Get(result.Raw, "container_port").Int(), p}
//			inarr = append(inarr, pm)
//		}
//		data["udp"] = inarr
//	}
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: data})
//}
