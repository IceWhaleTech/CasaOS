package v1

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
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
		c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
		return
	}
	name := c.PostForm("name")
	t := c.PostForm("type")
	language := c.GetHeader("Language")

	if len(name) == 0 || len(t) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	service.MyService.Casa().PushAppAnalyse(config.ServerInfo.Token, t, name, language)
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}
