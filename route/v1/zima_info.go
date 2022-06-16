/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-16 17:54:10
 * @FilePath: /CasaOS/route/v1/zima_info.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取信息系统信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/sysinfo [get]
func SysInfo(c *gin.Context) {
	info := service.MyService.ZiMa().GetSysInfo()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: info})
}
