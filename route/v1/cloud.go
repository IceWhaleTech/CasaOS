package v1

import (
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/dropbox"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/drivers/onedrive"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListStorages(c *gin.Context) {
	// var req model.PageReq
	// if err := c.ShouldBind(&req); err != nil {
	// 	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
	// 	return
	// }
	// req.Validate()

	//logger.Info("ListStorages", zap.Any("req", req))
	//storages, total, err := service.MyService.Storage().GetStorages(req.Page, req.PerPage)
	// if err != nil {
	// 	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	// 	return
	// }
	// c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: model.PageResp{
	// 	Content: storages,
	// 	Total:   total,
	// }})
	r, err := service.MyService.Storage().GetStorages()

	if err != nil {
		c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}

	for i := 0; i < len(r.MountPoints); i++ {
		dataMap, err := service.MyService.Storage().GetConfigByName(r.MountPoints[i].Fs)
		if err != nil {
			logger.Error("GetConfigByName", zap.Any("err", err))
			continue
		}
		if dataMap["type"] == "drive" {
			r.MountPoints[i].Icon = google_drive.ICONURL
		}
		if dataMap["type"] == "dropbox" {
			r.MountPoints[i].Icon = dropbox.ICONURL
		}
		if dataMap["type"] == "onedrive" {

			r.MountPoints[i].Icon = onedrive.ICONURL
		}
		r.MountPoints[i].Name = dataMap["username"]
	}
	list := []httper.MountPoint{}

	for _, v := range r.MountPoints {
		list = append(list, httper.MountPoint{
			Fs:         v.Fs,
			Icon:       v.Icon,
			MountPoint: v.MountPoint,
			Name:       v.Name,
		})
	}

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

func UmountStorage(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)
	mountPoint := json["mount_point"]
	if mountPoint == "" {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "mount_point is empty"})
		return
	}
	err := service.MyService.Storage().UnmountStorage(mountPoint)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	service.MyService.Storage().DeleteConfigByName(strings.ReplaceAll(mountPoint, "/mnt/", ""))
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: "success"})
}

func GetStorage(c *gin.Context) {

	// idStr := c.Query("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
	// 	return
	// }
	// storage, err := service.MyService.Storage().GetStorageById(uint(id))
	// if err != nil {
	// 	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	// 	return
	// }
	// c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: storage})
}
