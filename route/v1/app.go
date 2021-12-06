package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/docker"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/env_helper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	port2 "github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
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

	//service.MyService.Docker().DockerContainerCommit("test2")

	index := c.DefaultQuery("index", "1")
	size := c.DefaultQuery("size", "10")
	t := c.DefaultQuery("type", "rank")
	categoryId := c.DefaultQuery("category_id", "0")
	key := c.DefaultQuery("key", "")
	list, count := service.MyService.OAPI().GetServerList(index, size, t, categoryId, key)
	for i := 0; i < len(list); i++ {
		ct, _ := service.MyService.Docker().DockerListByImage(list[i].Image, list[i].ImageVersion)
		if ct != nil {
			list[i].State = ct.State
		}
	}
	data := make(map[string]interface{}, 2)
	data["count"] = count
	data["items"] = list

	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: data})
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
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: p})
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
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: port2.IsPortAvailable(p, t)})
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
// @Router /app/mylist [get]
func MyAppList(c *gin.Context) {
	index, _ := strconv.Atoi(c.DefaultQuery("index", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "0"))
	position, _ := strconv.ParseBool(c.DefaultQuery("position", "true"))
	list := service.MyService.App().GetMyList(index, size, position)
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
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
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
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
	info := service.MyService.OAPI().GetServerAppInfo(id)
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
	}

	for i := 0; i < len(info.Devices); i++ {
		if !file.CheckNotExist(info.Devices[i].ContainerPath) {
			info.Devices[i].Path = info.Devices[i].ContainerPath
		}
	}
	if len(info.Tip) > 0 {
		info.Tip = env_helper.ReplaceStringDefaultENV(info.Tip)
	}

	for i := 0; i < len(info.Volumes); i++ {
		info.Volumes[i].Path = docker.GetDir("", info.Volumes[i].ContainerPath)
	}
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

	//sort
	// if info.NetworkModel != "host" {
	// 	sort.PortsSort(portOrder).Sort(info.Configures.TcpPorts)
	// 	sort.PortsSort(portOrder).Sort(info.Configures.UdpPorts)
	// }

	// sort.EnvSort(envOrder).Sort(info.Envs)
	// sort.VolSort(volOrder).Sort(info.Volumes.([]model.PathMap))
	// sort.DevSort(devOrder).Sort(info.Devices)

	info.MaxMemory = service.MyService.ZiMa().GetMemInfo().Total >> 20

	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 获取远程分类列表
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/category [get]
func CategoryList(c *gin.Context) {
	list := service.MyService.OAPI().GetServerCategoryList()
	var count uint = 0
	for _, category := range list {
		count += category.Count
	}

	rear := append([]model.ServerCategoryList{}, list[0:]...)
	list = append(list[:0], model.ServerCategoryList{Count: count, Name: "All"})
	list = append(list, rear...)
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
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
	content := service.MyService.OAPI().ShareAppFile(str)
	c.JSON(http.StatusOK, json.RawMessage(content))
}

// @Summary Resource Usage
// @Produce  application/json
// @Accept application/json
// @Tags app
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /app/share [post]
func AppListResourceUsage() {

}
