package v1

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/gin-gonic/gin"
)

func SyncToSyncthing(c *gin.Context) {
	u := c.Param("url")
	target := "http://" + strings.Split(c.Request.Host, ":")[0] + ":" + config.SystemConfigInfo.SyncPort
	remote, err := url.Parse(target)
	if err != nil {
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.Header.Add("X-API-Key", config.SystemConfigInfo.SyncKey)
	//c.Request.Header.Add("X-API-Key", config.SystemConfigInfo.SyncKey)
	c.Request.URL.Path = u

	proxy.ServeHTTP(c.Writer, c.Request)
}

func GetSyncConfig(c *gin.Context) {
	data := make(map[string]string)
	data["key"] = config.SystemConfigInfo.SyncKey
	data["port"] = config.SystemConfigInfo.SyncPort
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: data})
}
