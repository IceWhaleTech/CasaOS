package v1

import (
	"bytes"
	json2 "encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/docker"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
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
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Duration(time.Second * 5),
}

// 打开docker的terminal
func DockerTerminal(c *gin.Context) {
	col := c.DefaultQuery("cols", "100")
	row := c.DefaultQuery("rows", "30")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	defer conn.Close()
	container := c.Param("id")
	hr, err := service.Exec(container, row, col)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
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

func PostSshLogin(c *gin.Context) {
	j := make(map[string]string)
	c.ShouldBind(&j)
	userName := j["username"]
	password := j["password"]
	port := j["port"]
	if userName == "" || password == "" || port == "" {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "Username or password or port is empty"})
		return
	}
	_, err := docker.NewSshClient(userName, password, port)
	if err != nil {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Please check if the username and port are correct, and make sure that ssh server is installed."})
		loger.Error("connect ssh error", zap.Any("error", err))
		return
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

func WsSsh(c *gin.Context) {
	_, e := exec.LookPath("ssh")
	if e != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: "ssh server not found"})
		return
	}

	userName := c.Query("username")
	password := c.Query("password")
	port := c.Query("port")
	wsConn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	var logBuff = new(bytes.Buffer)

	quitChan := make(chan bool, 3)
	// user := ""
	// password := ""
	var login int = 1
	cols, _ := strconv.Atoi(c.DefaultQuery("cols", "200"))
	rows, _ := strconv.Atoi(c.DefaultQuery("rows", "32"))
	var client *ssh.Client
	for login != 0 {

		var err error
		if userName == "" || password == "" || port == "" {
			wsConn.WriteMessage(websocket.TextMessage, []byte("username or password or port is empty"))
		}
		// wsConn.WriteMessage(websocket.TextMessage, []byte("login:"))
		// user = docker.ReceiveWsMsgUser(wsConn, logBuff)
		// wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		// wsConn.WriteMessage(websocket.TextMessage, []byte("password:"))
		// password = docker.ReceiveWsMsgPassword(wsConn, logBuff)
		// wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		client, err = docker.NewSshClient(userName, password, port)

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
// @Router /app/install [post]
func InstallApp(c *gin.Context) {
	var appInfo model.ServerAppList
	m := model.CustomizationPostData{}
	c.ShouldBind(&m)

	const CUSTOM = "custom"
	var dockerImage string
	var dockerImageVersion string

	//check app name is exist
	if len(m.Protocol) == 0 {
		m.Protocol = "http"
	}
	m.ContainerName = strings.Replace(m.Label, " ", "_", -1)
	if m.Origin != CUSTOM {
		oldName := m.ContainerName
		oldLabel := m.Label
		for i := 0; true; i++ {
			if i != 0 {
				m.ContainerName = oldName + "-" + strconv.Itoa(i)
				m.Label = oldLabel + "-" + strconv.Itoa(i)
			}
			if _, err := service.MyService.Docker().DockerListByName(m.ContainerName); err != nil {
				break
			}
		}
	} else {
		if _, err := service.MyService.Docker().DockerListByName(m.ContainerName); err == nil {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.ERROR_APP_NAME_EXIST, Message: common_err.GetMsg(common_err.ERROR_APP_NAME_EXIST)})
			return
		}

	}

	//check port
	if len(m.PortMap) > 0 && m.PortMap != "0" {
		//c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		portMap, _ := strconv.Atoi(m.PortMap)
		if !port2.IsPortAvailable(portMap, "tcp") {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + m.PortMap})
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
	m.Image = dockerImage + ":" + dockerImageVersion
	for _, u := range m.Ports {

		if u.Protocol == "udp" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "tcp" {

			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "both" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		}
	}
	if m.Origin == CUSTOM {
		for _, device := range m.Devices {
			if file.CheckNotExist(device.Path) {
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.DEVICE_NOT_EXIST, Message: device.Path + "," + common_err.GetMsg(common_err.DEVICE_NOT_EXIST)})
				return
			}

		}
	} else {
		dev := []model.PathMap{}
		for _, device := range dev {
			if !file.CheckNotExist(device.Path) {
				dev = append(dev, device)
			}
		}
		m.Devices = dev
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
	m.CustomId = id
	var relyMap = make(map[string]string)
	go func() {
		// installLog := model2.AppNotify{}
		// installLog.State = 0
		// installLog.CustomId = m.Label
		// installLog.Message = "installing rely"
		// installLog.Class = types.NOTIFY_APP
		// installLog.Type = types.NOTIFY_TYPE_UNIMPORTANT
		// installLog.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		// installLog.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		// installLog.Id = uuid.NewV4().String()
		// service.MyService.Notify().AddLog(installLog)
		if m.Origin != CUSTOM {
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
						rely.ContainerCustomId = m.Label
						var mysqlConfig model2.MysqlConfigs

						//结构体转换
						copier.Copy(&mysqlConfig, &mc)
						rely.Config = mysqlConfig
						service.MyService.Rely().Create(rely)

						relyMap["mysql"] = mid

					} else {
						docker_base.MysqlDelete(mysqlContainerId)
						// installLog.State = 0
						// installLog.Message = err.Error()
						// service.MyService.Notify().UpdateLog(installLog)
					}
				}
			}
		}

		// step：下载镜像
		err := service.MyService.Docker().DockerPullImage(dockerImage+":"+dockerImageVersion, m.Icon, m.Label)
		if err != nil {
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "PULLING"
			notify.Type = "INSTALL"
			notify.Success = false
			notify.Finished = false
			notify.Message = err.Error()
			service.MyService.Notify().SendInstallAppBySocket(notify)
			return
		}

		for !service.MyService.Docker().IsExistImage(m.Image) {
			time.Sleep(time.Second)
		}

		_, err = service.MyService.Docker().DockerContainerCreate(m, "")
		if err != nil {
			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":80}", 100)
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "STARTING"
			notify.Type = "INSTALL"
			notify.Success = false
			notify.Finished = false
			notify.Message = err.Error()
			service.MyService.Notify().SendInstallAppBySocket(notify)
			return
		} else {
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "STARTING"
			notify.Type = "INSTALL"
			notify.Success = true
			notify.Finished = false
			service.MyService.Notify().SendInstallAppBySocket(notify)
		}

		//		echo -e "hellow\nworld" >>

		//step：启动容器
		err = service.MyService.Docker().DockerContainerStart(m.ContainerName)
		if err != nil {
			//service.MyService.Redis().Set(id, "{\"id\"\""+id+"\",\"state\":false,\"message\":\""+err.Error()+"\",\"speed\":90}", 100)
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "STARTING"
			notify.Type = "INSTALL"
			notify.Success = false
			notify.Finished = false
			notify.Message = err.Error()
			service.MyService.Notify().SendInstallAppBySocket(notify)
			return
		} else {
			// if m.Origin != CUSTOM {
			// 	installLog.Message = "setting upnp"
			// } else {
			// 	installLog.Message = "nearing completion"
			// }
			// service.MyService.Notify().UpdateLog(installLog)
		}

		//step: 启动成功     检查容器状态确认启动成功
		container, err := service.MyService.Docker().DockerContainerInfo(m.ContainerName)
		if err != nil && container.ContainerJSONBase.State.Running {
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "INSTALLED"
			notify.Type = "INSTALL"
			notify.Success = false
			notify.Finished = true
			notify.Message = err.Error()
			service.MyService.Notify().SendInstallAppBySocket(notify)
			return
		} else {
			notify := notify.Application{}
			notify.Icon = m.Icon
			notify.Name = m.Label
			notify.State = "INSTALLED"
			notify.Type = "INSTALL"
			notify.Success = true
			notify.Finished = true
			service.MyService.Notify().SendInstallAppBySocket(notify)
		}

		// if m.Origin != "custom" {
		// 	for i := 0; i < len(m.Volumes); i++ {
		// 		m.Volumes[i].Path = docker.GetDir(id, m.Volumes[i].Path)
		// 	}
		// }
		//service.MyService.App().SaveContainer(md)
		config.CasaOSGlobalVariables.AppChange = true

	}()

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: m.Label})

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
//	c.ShouldBind(&m)
//	//检查端口
//	if len(m.PortMap) == 0 || m.PortMap == "0" {
//		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
//		return
//	}
//	if len(m.Port) == 0 || m.Port == "0" {
//		m.Port = m.PortMap
//	}
//
//	portMap, _ := strconv.Atoi(m.PortMap)
//	if !port2.IsPortAvailable(portMap, "tcp") {
//		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: "Duplicate port:" + m.PortMap})
//		return
//	}
//
//	for _, u := range m.Udp {
//		t, _ := strconv.Atoi(u.CommendPort)
//		if !port2.IsPortAvailable(t, "udp") {
//			c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: "Duplicate port:" + u.CommendPort})
//			return
//		}
//	}
//
//	for _, t := range m.Tcp {
//		te, _ := strconv.Atoi(t.CommendPort)
//		if !port2.IsPortAvailable(te, "tcp") {
//			c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: "Duplicate port:" + t.CommendPort})
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
//		c.JSON(http.StatusOK, model.Result{Success: common_err.PULL_IMAGE_ERROR, Message: common_err.GetMsg(common_err.PULL_IMAGE_ERROR)})
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
//	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: id})
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
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	//info := service.MyService.App().GetUninstallInfo(appId)

	info, err := service.MyService.Docker().DockerContainerInfo(appId)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}

	//step：停止容器
	err = service.MyService.Docker().DockerContainerStop(appId)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.UNINSTALL_APP_ERROR, Message: common_err.GetMsg(common_err.UNINSTALL_APP_ERROR), Data: err.Error()})
		return
	}

	err = service.MyService.Docker().DockerContainerRemove(appId, false)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.UNINSTALL_APP_ERROR, Message: common_err.GetMsg(common_err.UNINSTALL_APP_ERROR), Data: err.Error()})
		return
	}

	// step：remove image
	service.MyService.Docker().DockerImageRemove(info.Config.Image)

	if info.Config.Labels["origin"] != "custom" {
		//step: 删除文件夹
		for _, v := range info.Mounts {
			if strings.Contains(v.Source, info.Name) {
				path := filepath.Join(strings.Split(v.Source, info.Name)[0], info.Name)
				service.MyService.App().DelAppConfigDir(path)
			}
		}

		//step: 删除install log
		//service.MyService.Notify().DelLog(appId)

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
	config.CasaOSGlobalVariables.AppChange = true
	notify := notify.Application{}
	notify.Icon = info.Config.Labels["icon"]
	notify.Name = strings.ReplaceAll(info.Name, "/", "")
	notify.State = "FINISHED"
	notify.Type = "UNINSTALL"
	notify.Success = true
	notify.Finished = true
	service.MyService.Notify().SendUninstallAppBySocket(notify)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})

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
	js := make(map[string]string)
	c.ShouldBind(&js)
	state := js["state"]
	if len(appId) == 0 || len(state) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	var err error
	if state == "start" {
		err = service.MyService.Docker().DockerContainerStart(appId)
	} else if state == "restart" {
		service.MyService.Docker().DockerContainerStop(appId)
		err = service.MyService.Docker().DockerContainerStart(appId)
	} else {
		err = service.MyService.Docker().DockerContainerStop(appId)
	}

	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	info, err := service.MyService.App().GetContainerInfo(appId)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}

	// @tiger - 用 {'state': ...} 来体现出参上下文
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: info.State})
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
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: string(log)})
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
	//t := c.DefaultQuery("type", "0")
	containerInfo, e := service.MyService.App().GetSimpleContainerInfo(id)
	if e != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: e.Error()})
		return
	}

	var data = make(map[string]interface{})

	data["state"] = containerInfo.State

	// if t == "1" {
	// 	appInfo := service.MyService.App().GetAppDBInfo(id)
	// 	data["app"] = appInfo
	// }

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
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
	c.ShouldBind(&m)

	if len(id) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	//var cpd model.CustomizationPostData

	//copier.Copy(&cpd, &m)

	//appInfo := service.MyService.App().GetAppDBInfo(id)
	//info, err := service.MyService.Docker().DockerContainerInfo(id)

	// //check app name is exist
	// if _, err := service.MyService.Docker().DockerListByName(m.Label); err == nil {
	// 	c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR_APP_NAME_EXIST, Message: common_err.GetMsg(common_err.ERROR_APP_NAME_EXIST)})
	// 	return
	// }
	service.MyService.Docker().DockerContainerStop(id)
	portMap, _ := strconv.Atoi(m.PortMap)
	if !port2.IsPortAvailable(portMap, "tcp") {
		service.MyService.Docker().DockerContainerStart(id)
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + m.PortMap})
		return
	}

	for _, u := range m.Ports {

		if u.Protocol == "udp" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				service.MyService.Docker().DockerContainerStart(id)
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "tcp" {
			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				service.MyService.Docker().DockerContainerStart(id)
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		} else if u.Protocol == "both" {
			t, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(t, "udp") {
				service.MyService.Docker().DockerContainerStart(id)
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}

			te, _ := strconv.Atoi(u.CommendPort)
			if !port2.IsPortAvailable(te, "tcp") {
				service.MyService.Docker().DockerContainerStart(id)
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: "Duplicate port:" + u.CommendPort})
				return
			}
		}

	}
	service.MyService.Docker().DockerContainerUpdateName(id, id)
	//service.MyService.Docker().DockerContainerRemove(id, true)

	containerId, err := service.MyService.Docker().DockerContainerCreate(m, id)
	if err != nil {
		service.MyService.Docker().DockerContainerUpdateName(m.ContainerName, id)
		service.MyService.Docker().DockerContainerStart(id)
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	//		echo -e "hellow\nworld" >>

	//step：启动容器
	err = service.MyService.Docker().DockerContainerStart(containerId)

	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR)})
		return
	}
	service.MyService.Docker().DockerContainerRemove(id, true)
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

	//service.MyService.App().UpdateApp(appInfo)

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary update app version
// @Produce  application/json
// @Accept multipart/form-data
// @Tags app
// @Param  id path string true "容器id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/update/{id} [put]
func PutAppUpdate(c *gin.Context) {
	id := c.Param("id")

	if len(id) == 0 {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}

	inspect, err := service.MyService.Docker().DockerContainerInfo(id)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return

	}
	imageLatest := strings.Split(inspect.Config.Image, ":")[0] + ":latest"
	err = service.MyService.Docker().DockerPullImage(imageLatest, "", "")
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return

	}
	service.MyService.Docker().DockerContainerStop(id)
	service.MyService.Docker().DockerContainerUpdateName(id, id)
	//service.MyService.Docker().DockerContainerRemove(id, true)
	inspect.Image = imageLatest
	inspect.Config.Image = imageLatest
	containerId, err := service.MyService.Docker().DockerContainerCopyCreate(inspect)
	if err != nil {
		service.MyService.Docker().DockerContainerUpdateName(inspect.Name, id)
		service.MyService.Docker().DockerContainerStart(id)
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR)})
		return
	}

	//step：启动容器
	err = service.MyService.Docker().DockerContainerStart(containerId)

	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR)})
		return
	}
	service.MyService.Docker().DockerContainerRemove(id, true)
	delete(service.NewVersionApp, id)

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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

	// @tiger - 作为最佳实践，不应该直接把数据库的信息返回，来避免未来数据库结构上的迭代带来的新字段
	appInfo := service.MyService.App().GetAppDBInfo(appId)
	containerInfo, _ := service.MyService.Docker().DockerContainerStats(appId)
	var cpuModel = "arm"
	if cpu := service.MyService.System().GetCpuInfo(); len(cpu) > 0 {
		if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "intel") > 0 {
			cpuModel = "intel"
		} else if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "amd") > 0 {
			cpuModel = "amd"
		}
	}

	info, err := service.MyService.Docker().DockerContainerInfo(appId)
	if err != nil {
		//todo 需要自定义错误
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	con := struct {
		Status    string `json:"status"`
		StartedAt string `json:"started_at"`
		CPUShares int64  `json:"cpu_shares"`
		Memory    int64  `json:"total_memory"`   // @tiger - 改成 total_memory，方便以后增加 free_memory 之类的字段
		Restart   string `json:"restart_policy"` // @tiger - 改成 restart_policy?
	}{Status: info.State.Status, StartedAt: info.State.StartedAt, CPUShares: info.HostConfig.CPUShares, Memory: info.HostConfig.Memory >> 20, Restart: info.HostConfig.RestartPolicy.Name}
	data := make(map[string]interface{}, 5)
	data["app"] = appInfo                                             // @tiget - 最佳实践是，返回 appid，然后具体的 app 信息由前端另行获取
	data["cpu"] = cpuModel                                            // @tiger - 改成 arch
	data["memory"] = service.MyService.System().GetMemInfo()["total"] // @tiger - 改成 total_memory，方便以后增加 free_memory 之类的字段
	data["container"] = json2.RawMessage(containerInfo)
	data["info"] = con
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

func GetDockerNetworks(c *gin.Context) {
	networks := service.MyService.Docker().DockerNetworkModelList()
	list := []map[string]string{}
	for _, network := range networks {
		if network.Driver != "null" {
			list = append(list, map[string]string{"name": network.Name, "driver": network.Driver, "id": network.ID})
		}
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
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
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: appInfo})
}

// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path string true "appid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/update/{id}/info [get]
func ContainerUpdateInfo(c *gin.Context) {
	appId := c.Param("id")
	//appInfo := service.MyService.App().GetAppDBInfo(appId)
	info, err := service.MyService.Docker().DockerContainerInfo(appId)
	if err != nil {

		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: err.Error()})
		return
	}
	var port model.PortArray
	// json2.Unmarshal([]byte(appInfo.Ports), &port)

	for k, v := range info.HostConfig.PortBindings {
		temp := model.PortMap{
			CommendPort:   v[0].HostPort,
			ContainerPort: k.Port(),

			Protocol: k.Proto(),
		}
		port = append(port, temp)
	}

	var envs model.EnvArray
	// json2.Unmarshal([]byte(appInfo.Envs), &envs)

	showENV := info.Config.Labels["show_env"]
	showENVList := strings.Split(showENV, ",")
	showENVMap := make(map[string]string)
	if len(showENVList) > 0 && showENVList[0] != "" {
		for _, name := range showENVList {
			showENVMap[name] = "1"
		}
	}
	for _, v := range info.Config.Env {
		if len(showENVList) > 0 {
			if _, ok := showENVMap[strings.Split(v, "=")[0]]; ok {
				temp := model.Env{
					Name:  strings.Split(v, "=")[0],
					Value: strings.Split(v, "=")[1],
				}
				envs = append(envs, temp)
			}
		} else {
			temp := model.Env{
				Name:  strings.Split(v, "=")[0],
				Value: strings.Split(v, "=")[1],
			}
			envs = append(envs, temp)
		}

	}

	var vol model.PathArray
	// json2.Unmarshal([]byte(appInfo.Volumes), &vol)

	for i := 0; i < len(info.Mounts); i++ {
		temp := model.PathMap{
			Path:          strings.ReplaceAll(info.Mounts[i].Source, "$AppID", info.Name),
			ContainerPath: info.Mounts[i].Destination,
		}
		vol = append(vol, temp)
	}
	var driver model.PathArray

	//volumesStr, _ := json2.Marshal(m.Volumes)
	//devicesStr, _ := json2.Marshal(m.Devices)
	for _, v := range info.HostConfig.Resources.Devices {
		temp := model.PathMap{
			Path:          v.PathOnHost,
			ContainerPath: v.PathInContainer,
		}
		driver = append(driver, temp)
	}

	m := model.CustomizationPostData{}
	m.Icon = info.Config.Labels["icon"]
	m.Ports = port
	m.Image = info.Config.Image
	m.Origin = info.Config.Labels["origin"]
	if len(m.Origin) == 0 {
		m.Origin = "local"
	}
	m.NetworkModel = string(info.HostConfig.NetworkMode)
	m.Description = info.Config.Labels["desc"]
	m.ContainerName = strings.ReplaceAll(info.Name, "/", "")
	m.PortMap = info.Config.Labels["web"]
	m.Devices = driver
	m.Envs = envs
	m.Memory = info.HostConfig.Memory >> 20
	m.CpuShares = info.HostConfig.CPUShares
	m.Volumes = vol //appInfo.Volumes
	m.Restart = info.HostConfig.RestartPolicy.Name
	m.EnableUPNP = false
	m.Index = info.Config.Labels["index"]
	m.Position = false
	m.CustomId = info.Config.Labels["custom_id"]
	m.Host = info.Config.Labels["host"]
	if len(m.CustomId) == 0 {
		m.CustomId = uuid.NewV4().String()
	}
	m.CapAdd = info.HostConfig.CapAdd
	m.Cmd = info.Config.Cmd
	m.HostName = info.Config.Hostname
	m.Privileged = info.HostConfig.Privileged
	name := info.Config.Labels["name"]
	if len(name) == 0 {
		name = strings.ReplaceAll(info.Name, "/", "")
	}
	m.Label = name

	m.Protocol = info.Config.Labels["protocol"]
	if m.Protocol == "" {
		m.Protocol = "http"
	}

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: m})
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
//	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
//}
