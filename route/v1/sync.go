/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-11-08 18:02:02
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-21 19:13:59
 * @FilePath: /CasaOS/route/v1/sync.go
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
	"github.com/gin-gonic/gin"
)

func GetSyncConfig(c *gin.Context) {
	data := make(map[string]string)
	data["key"] = config.SystemConfigInfo.SyncKey
	data["port"] = config.SystemConfigInfo.SyncPort
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}
