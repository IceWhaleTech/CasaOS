package v1

import (
	"net/http/httputil"
	"net/url"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/gin-gonic/gin"
)

func SyncToSyncthing(c *gin.Context) {
	u := c.Param("url")
	target := "http://127.0.0.1:" + config.SystemConfigInfo.SyncPort
	remote, err := url.Parse(target)
	if err != nil {
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.Header.Add("X-API-Key", config.SystemConfigInfo.SyncKey)
	c.Request.URL.Path = u
	proxy.ServeHTTP(c.Writer, c.Request)
}
