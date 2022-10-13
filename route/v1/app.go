package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"

	port2 "github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

const (
	dockerRootDirFilePath             = "/var/lib/casaos/docker_root"
	dockerDaemonConfigurationFilePath = "/etc/docker/daemon.json"
)

// @Summary 获取远程列表
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param index query int false "页码"
// @Param size query int false "每页数量"
// @Param  category_id query int false "分类id"
// @Param  type query string false "rank,new"
// @Param  key query string false "search key"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/list [get]
func AppList(c *gin.Context) {
	// service.MyService.Docker().DockerContainerCommit("test2")

	index := c.DefaultQuery("index", "1")
	size := c.DefaultQuery("size", "10000")
	t := c.DefaultQuery("type", "rank")
	categoryId := c.DefaultQuery("category_id", "0")
	key := c.DefaultQuery("key", "")
	if len(index) == 0 || len(size) == 0 || len(t) == 0 || len(categoryId) == 0 {
		c.JSON(common_err.CLIENT_ERROR, &model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	collection, err := service.MyService.Casa().GetServerList(index, size, t, categoryId, key)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, &model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	// for i := 0; i < len(recommend); i++ {
	// 	ct, _ := service.MyService.Docker().DockerListByImage(recommend[i].Image, recommend[i].ImageVersion)
	// 	if ct != nil {
	// 		recommend[i].State = ct.State
	// 	}
	// }
	// for i := 0; i < len(list); i++ {
	// 	ct, _ := service.MyService.Docker().DockerListByImage(list[i].Image, list[i].ImageVersion)
	// 	if ct != nil {
	// 		list[i].State = ct.State
	// 	}
	// }
	// for i := 0; i < len(community); i++ {
	// 	ct, _ := service.MyService.Docker().DockerListByImage(community[i].Image, community[i].ImageVersion)
	// 	if ct != nil {
	// 		community[i].State = ct.State
	// 	}
	// }
	data := make(map[string]interface{}, 3)
	data["recommend"] = collection.Recommend
	data["list"] = collection.List
	data["community"] = collection.Community

	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

// @Summary 获取一个可用端口
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  type query string true "端口类型 udp/tcp"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/getport [get]
func GetPort(c *gin.Context) {
	t := c.DefaultQuery("type", "tcp")
	var p int
	ok := true
	for ok {
		p, _ = port2.GetAvailablePort(t)
		ok = !port2.IsPortAvailable(p, t)
	}
	// @tiger 这里最好封装成 {'port': ...} 的形式，来体现出参的上下文
	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: p})
}

// @Summary 检查端口是否可用
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  port path int true "端口号"
// @Param  type query string true "端口类型 udp/tcp"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/check/{port} [get]
func PortCheck(c *gin.Context) {
	p, _ := strconv.Atoi(c.Param("port"))
	t := c.DefaultQuery("type", "tcp")
	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: port2.IsPortAvailable(p, t)})
}

// @Summary 我的应用列表
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Param  index query int false "index"
// @Param  size query int false "size"
// @Param  position query bool false "是否是首页应用"
// @Success 200 {string} string "ok"
// @Router /app/my/list [get]
func MyAppList(c *gin.Context) {
	index, _ := strconv.Atoi(c.DefaultQuery("index", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "0"))
	position, _ := strconv.ParseBool(c.DefaultQuery("position", "true"))
	list, unTranslation := service.MyService.App().GetMyList(index, size, position)
	data := make(map[string]interface{}, 2)
	data["casaos_apps"] = list
	data["local_apps"] = unTranslation

	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

// @Summary my app hardware usage list
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/usage [get]
func AppUsageList(c *gin.Context) {
	list := service.MyService.App().GetHardwareUsage()
	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
	// c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: nil})
}

// @Summary 应用详情
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Param  id path int true "id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/appinfo/{id} [get]
func AppInfo(c *gin.Context) {
	id := c.Param("id")
	language := c.GetHeader("Language")
	info, err := service.MyService.Casa().GetServerAppInfo(id, "", language)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, &model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	if info.NetworkModel != "host" {
		for i := 0; i < len(info.Ports); i++ {
			if p, _ := strconv.Atoi(info.Ports[i].ContainerPort); port2.IsPortAvailable(p, info.Ports[i].Protocol) {
				info.Ports[i].CommendPort = strconv.Itoa(p)
			} else {
				if info.Ports[i].Protocol == "tcp" {
					if p, err := port2.GetAvailablePort("tcp"); err == nil {
						info.Ports[i].CommendPort = strconv.Itoa(p)
					}
				} else if info.Ports[i].Protocol == "upd" {
					if p, err := port2.GetAvailablePort("udp"); err == nil {
						info.Ports[i].CommendPort = strconv.Itoa(p)
					}
				}
			}

			if info.Ports[i].Type == 0 {
				info.PortMap = info.Ports[i].CommendPort
			}
		}
	} else {
		for i := 0; i < len(info.Ports); i++ {
			if info.Ports[i].Type == 0 {
				info.PortMap = info.Ports[i].ContainerPort
				break
			}
		}
	}

	for i := 0; i < len(info.Devices); i++ {
		if !file.CheckNotExist(info.Devices[i].ContainerPath) {
			info.Devices[i].Path = info.Devices[i].ContainerPath
		}
	}
	// if len(info.Tip) > 0 {
	// 	info.Tip = env_helper.ReplaceStringDefaultENV(info.Tip)
	// }

	// portOrder := func(c1, c2 *model.Ports) bool {
	// 	return c1.Type < c2.Type
	// }

	// envOrder := func(c1, c2 *model.Envs) bool {
	// 	return c1.Type < c2.Type
	// }

	// volOrder := func(c1, c2 *model.Volume) bool {
	// 	return c1.Type < c2.Type
	// }

	// devOrder := func(c1, c2 *model.Devices) bool {
	// 	return c1.Type < c2.Type
	// }

	// sort
	// if info.NetworkModel != "host" {
	// 	sort.PortsSort(portOrder).Sort(info.Configures.TcpPorts)
	// 	sort.PortsSort(portOrder).Sort(info.Configures.UdpPorts)
	// }

	// sort.EnvSort(envOrder).Sort(info.Envs)
	// sort.VolSort(volOrder).Sort(info.Volumes.([]model.PathMap))
	// sort.DevSort(devOrder).Sort(info.Devices)
	info.Image += ":" + info.ImageVersion
	info.MaxMemory = (service.MyService.System().GetMemInfo()["total"]).(uint64) >> 20

	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: info})
}

// @Summary 获取远程分类列表
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/category [get]
func CategoryList(c *gin.Context) {
	list, err := service.MyService.Casa().GetServerCategoryList()
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, &model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	var count uint = 0
	for _, category := range list {
		count += category.Count
	}

	rear := append([]model.CategoryList{}, list[0:]...)
	list = append(list[:0], model.CategoryList{Count: count, Name: "All", Font: "apps"})
	list = append(list, rear...)
	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

// @Summary 分享该应用配置
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/share [post]
func ShareAppFile(c *gin.Context) {
	str, _ := ioutil.ReadAll(c.Request.Body)
	content := service.MyService.Casa().ShareAppFile(str)
	c.JSON(common_err.SUCCESS, json.RawMessage(content))
}

func GetDockerDaemonConfiguration(c *gin.Context) {
	// info, err := service.MyService.Docker().GetDockerInfo()
	// if err != nil {
	// 	c.JSON(common_err.SERVICE_ERROR, &model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	// 	return
	// }
	data := make(map[string]interface{})

	if file.Exists(dockerRootDirFilePath) {
		buf := file.ReadFullFile(dockerRootDirFilePath)
		err := json.Unmarshal(buf, &data)
		if err != nil {
			c.JSON(common_err.CLIENT_ERROR, &model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: err})
			return
		}
	}
	c.JSON(common_err.SUCCESS, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

func PutDockerDaemonConfiguration(c *gin.Context) {
	request := make(map[string]interface{})
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: err})
		return
	}

	value, ok := request["docker_root_dir"]
	if !ok {
		c.JSON(http.StatusBadRequest, &model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "`docker_root_dir` should not empty"})
		return
	}

	dockerConfig := model.DockerDaemonConfigurationModel{}
	if file.Exists(dockerDaemonConfigurationFilePath) {
		byteResult := file.ReadFullFile(dockerDaemonConfigurationFilePath)
		err := json.Unmarshal(byteResult, &dockerConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &model.Result{Success: common_err.SERVICE_ERROR, Message: "error when trying to deserialize " + dockerDaemonConfigurationFilePath, Data: err})
			return
		}
	}

	dockerRootDir := value.(string)
	if dockerRootDir == "/" {
		dockerConfig.Root = "" // omitempty - empty string will not be serialized
	} else {
		if !file.Exists(dockerRootDir) {
			c.JSON(http.StatusBadRequest, &model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.DIR_NOT_EXISTS), Data: common_err.GetMsg(common_err.DIR_NOT_EXISTS)})
			return
		}

		dockerConfig.Root = filepath.Join(dockerRootDir, "docker")

		if err := file.IsNotExistMkDir(dockerConfig.Root); err != nil {
			c.JSON(http.StatusInternalServerError, &model.Result{Success: common_err.SERVICE_ERROR, Message: "error when trying to create " + dockerConfig.Root, Data: err})
			return
		}
	}

	if buf, err := json.Marshal(request); err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{Success: common_err.CLIENT_ERROR, Message: "error when trying to serialize docker root json", Data: err})
		return
	} else {
		if err := file.WriteToFullPath(buf, dockerRootDirFilePath, 0o644); err != nil {
			c.JSON(http.StatusInternalServerError, &model.Result{Success: common_err.SERVICE_ERROR, Message: "error when trying to write " + dockerRootDirFilePath, Data: err})
			return
		}
	}

	if buf, err := json.Marshal(dockerConfig); err != nil {
		c.JSON(http.StatusBadRequest, &model.Result{Success: common_err.CLIENT_ERROR, Message: "error when trying to serialize docker config", Data: dockerConfig})
		return
	} else {
		if err := file.WriteToFullPath(buf, dockerDaemonConfigurationFilePath, 0o644); err != nil {
			c.JSON(http.StatusInternalServerError, &model.Result{Success: common_err.SERVICE_ERROR, Message: "error when trying to write to " + dockerDaemonConfigurationFilePath, Data: err})
			return
		}
	}

	println(command.ExecResultStr("systemctl daemon-reload"))
	println(command.ExecResultStr("systemctl restart docker"))

	c.JSON(http.StatusOK, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: request})
}
