/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-11-08 18:02:02
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 14:30:24
 * @FilePath: /CasaOS/route/v1/sync.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
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
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}
