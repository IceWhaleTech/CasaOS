/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-03-08 11:28:15
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 14:20:21
 * @FilePath: /CasaOS/route/v1/analyse.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary post app analyse
// @Produce  application/json
// @Accept multipart/form-data
// @Tags analyse
// @Param name formData string true "app name"
// @Param type formData string true "action" Enums(open,delete)
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /analyse/app [post]
func PostAppAnalyse(c *gin.Context) {
	if config.SystemConfigInfo.Analyse == "False" {
		c.JSON(http.StatusOK, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
		return
	}
	name := c.PostForm("name")
	t := c.PostForm("type")
	language := c.GetHeader("Language")

	if len(name) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}

	service.MyService.Casa().PushAppAnalyse(config.ServerInfo.Token, t, name, language)
	c.JSON(http.StatusOK, &model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
