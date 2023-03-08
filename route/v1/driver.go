package v1

import (
	"github.com/IceWhaleTech/CasaOS-Common/model"
	"github.com/IceWhaleTech/CasaOS-Common/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/internal/op"
	"github.com/gin-gonic/gin"
)

func ListDriverInfo(c *gin.Context) {
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: op.GetDriverInfoMap()})
}
