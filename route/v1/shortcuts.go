/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 14:30:05
 * @FilePath: /CasaOS/route/v1/shortcuts.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"
	"net/url"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/gin-gonic/gin"
)

// @Summary 获取短链列表
// @Produce  application/json
// @Accept application/json
// @Tags shortcuts
// @Param username formData string true "User name"
// @Param pwd  formData string true "password"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /shortcuts/list [get]
func GetShortcutsList(c *gin.Context) {
	list := service.MyService.Shortcuts().GetList()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

// @Summary 添加shortcuts
// @Produce  application/json
// @Accept application/json
// @Tags shortcuts
// @Param title formData string true "title"
// @Param url  formData string true "url"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /shortcuts/add [post]
func PostShortcutsAdd(c *gin.Context) {
	var m model2.ShortcutsDBModel

	c.BindJSON(&m)
	if len(m.Url) == 0 || len(m.Title) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	u, err := url.Parse(m.Url)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.SHORTCUTS_URL_ERROR, Message: common_err.GetMsg(common_err.SHORTCUTS_URL_ERROR), Data: err.Error()})
		return
	}
	m.Icon = "https://api.faviconkit.com/" + u.Host + "/57"
	service.MyService.Shortcuts().AddData(m)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})

}

// @Summary 删除shortcuts
// @Produce  application/json
// @Accept application/json
// @Tags shortcuts
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /shortcuts/del/{id} [post]
func DeleteShortcutsDelete(c *gin.Context) {
	id := c.Param("id")
	service.MyService.Shortcuts().DeleteData(id)
	c.JSON(http.StatusOK, model.Result{
		Success: common_err.SUCCESS,
		Message: common_err.GetMsg(common_err.SUCCESS),
		Data:    "",
	})
}

// @Summary 编辑shortcuts
// @Produce  application/json
// @Accept application/json
// @Tags shortcuts
// @Param title formData string true "title"
// @Param url  formData string true "url"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /shortcuts/edit [put]
func PutShortcutsEdit(c *gin.Context) {
	var m model2.ShortcutsDBModel
	c.BindJSON(&m)
	if len(m.Url) == 0 || len(m.Title) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	u, err := url.Parse(m.Url)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.SHORTCUTS_URL_ERROR, Message: common_err.GetMsg(common_err.SHORTCUTS_URL_ERROR), Data: err.Error()})
		return
	}
	m.Icon = "https://api.faviconkit.com/" + u.Host + "/57"
	service.MyService.Shortcuts().EditData(m)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: ""})
}
