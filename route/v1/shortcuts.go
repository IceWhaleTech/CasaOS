package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"oasis/model"
	"oasis/pkg/utils/oasis_err"
	"oasis/service"
	model2 "oasis/service/model"
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
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: list})
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
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	u, err := url.Parse(m.Url)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.SHORTCUTS_URL_ERROR, Message: oasis_err.GetMsg(oasis_err.SHORTCUTS_URL_ERROR), Data: err.Error()})
		return
	}
	m.Icon = "https://api.faviconkit.com/" + u.Host + "/57"
	service.MyService.Shortcuts().AddData(m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})

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
		Success: oasis_err.SUCCESS,
		Message: oasis_err.GetMsg(oasis_err.SUCCESS),
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
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	u, err := url.Parse(m.Url)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.SHORTCUTS_URL_ERROR, Message: oasis_err.GetMsg(oasis_err.SHORTCUTS_URL_ERROR), Data: err.Error()})
		return
	}
	m.Icon = "https://api.faviconkit.com/" + u.Host + "/57"
	service.MyService.Shortcuts().EditData(m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: ""})
}
