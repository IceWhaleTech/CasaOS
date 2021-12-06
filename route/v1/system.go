package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	port2 "github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/version"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
)

// @Summary 系统信息
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/chackversion [get]
func CheckVersion(c *gin.Context) {
	need, version := version.IsNeedUpdate()
	if need {
		installLog := model2.AppNotify{}
		installLog.State = 0
		installLog.Message = "New version " + version.Version + " is ready, ready to upgrade"
		installLog.Type = types.NOTIFY_TYPE_NEED_CONFIRM
		installLog.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.Name = "CasaOS System"
		service.MyService.Notify().AddLog(installLog)
	}
	data := make(map[string]interface{}, 1)
	data["is_need"] = need
	data["version"] = version
	data["current_version"] = types.CURRENTVERSION
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: data})
	return
}

// @Summary 系统信息
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/update [post]
func SystemUpdate(c *gin.Context) {
	need, version := version.IsNeedUpdate()
	if need {
		service.MyService.System().UpdateSystemVersion(version.Version)
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

//系统配置
func GetSystemConfig(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: json.RawMessage(config.SystemConfigInfo.ConfigStr)})
}

// @Summary  get logs
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/error/logs [get]
func GetCasaOSErrorLogs(c *gin.Context) {
	line, _ := strconv.Atoi(c.DefaultQuery("line", "100"))
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: service.MyService.System().GetCasaOSLogs(line)})
}

// @Summary 修改配置文件
// @Produce  application/json
// @Accept multipart/form-data
// @Tags sys
// @Param config formData string true "config json string"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/changhead [post]
func PostSetSystemConfig(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	service.MyService.System().UpSystemConfig(string(buf[0:n]), "")
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err.SUCCESS,
			Message: oasis_err.GetMsg(oasis_err.SUCCESS),
			Data:    json.RawMessage(config.SystemConfigInfo.ConfigStr),
		})
}

//系统配置
func GetSystemConfigDebug(c *gin.Context) {

	array := service.MyService.System().GetSystemConfigDebug()
	disk := service.MyService.ZiMa().GetDiskInfo()
	array = append(array, fmt.Sprintf("disk,total:%v,used:%v,UsedPercent:%v", disk.Total>>20, disk.Used>>20, disk.UsedPercent))

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: array})
}
func Sys(c *gin.Context) {
	service.DockerPull()
}

//widget配置
func GetWidgetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: json.RawMessage(config.SystemConfigInfo.WidgetList)})
}

// @Summary 修改组件配置文件
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/widget/config [post]
func PostSetWidgetConfig(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	service.MyService.System().UpSystemConfig("", string(buf[0:n]))
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err.SUCCESS,
			Message: oasis_err.GetMsg(oasis_err.SUCCESS),
			Data:    json.RawMessage(config.SystemConfigInfo.WidgetList),
		})
}

// @Summary edit casaos server port
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Param port formData file true "用户头像"
// @Success 200 {string} string "ok"
// @Router /sys/widget/config [post]
func PutCasaOSPort(c *gin.Context) {
	port, err := strconv.Atoi(c.PostForm("port"))
	if err != nil {
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err.ERROR,
				Message: err.Error(),
			})
		return
	}

	isAvailable := port2.IsPortAvailable(port, "tcp")
	if !isAvailable {
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err.PORT_IS_OCCUPIED,
				Message: oasis_err.GetMsg(oasis_err.PORT_IS_OCCUPIED),
			})
		return
	}
	service.MyService.System().UpSystemPort(strconv.Itoa(port))
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err.SUCCESS,
			Message: oasis_err.GetMsg(oasis_err.SUCCESS),
		})
}

// @Summary 检查是否进入引导状态
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /guide/check [get]
func GetGuideCheck(c *gin.Context) {
	initUser := false
	if !config.UserInfo.Initialized {
		initUser = true
	}
	data := make(map[string]interface{}, 1)
	data["need_init_user"] = initUser
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err.SUCCESS,
			Message: oasis_err.GetMsg(oasis_err.SUCCESS),
			Data:    data,
		})
}

// @Summary active killing casaos
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/kill [post]
func PostKillCasaOS(c *gin.Context) {
	os.Exit(0)
}
