package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"oasis/model"
	"oasis/pkg/config"
	"oasis/pkg/utils/oasis_err"
	"oasis/pkg/utils/version"
	"oasis/service"
	model2 "oasis/service/model"
	"oasis/types"
	"strconv"
	"time"
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
		installLog.CustomId = ""
		installLog.State = 0
		installLog.Message = "New version " + version.Version + " is ready, ready to upgrade"
		installLog.Speed = 100
		installLog.Type = types.NOTIFY_TYPE_NEED_CONFIRM
		installLog.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
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
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: config.SystemConfigInfo})
}

// @Summary 修改配置文件
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param file formData file true "用户头像"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/changhead [post]
func PostSetSystemConfig(c *gin.Context) {
	var systemConfig model.SystemConfig
	c.BindJSON(&systemConfig)
	service.MyService.System().UpSystemConfig(systemConfig)
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err.SUCCESS,
			Message: oasis_err.GetMsg(oasis_err.SUCCESS),
			Data:    config.SystemConfigInfo,
		})
	return
}

//系统配置
func GetSystemConfigDebug(c *gin.Context) {

	array := service.MyService.System().GetSystemConfigDebug()
	disk := service.MyService.ZiMa().GetDiskInfo()
	array = append(array, fmt.Sprintf("disk,totle:%v,used:%v,UsedPercent:%v", disk.Total>>20, disk.Used>>20, disk.UsedPercent))

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: array})
}
func Sys(c *gin.Context) {
	service.DockerPull()
}
