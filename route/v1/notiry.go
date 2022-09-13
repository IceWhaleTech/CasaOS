package v1

import (
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

func PostNotifyMssage(c *gin.Context) {
	path := c.Param("path")
	message := make(map[string]interface{})
	c.ShouldBind(&message)
	service.MyService.Notify().SendNotify(path, message)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
func PostSystemNotyfiy(c *gin.Context) {
	message := make(map[string]interface{})
	c.ShouldBind(&message)
	service.MyService.Notify().SettingSystemTempData(message)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
