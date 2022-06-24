/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-23 17:27:43
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-23 17:27:48
 * @FilePath: /CasaOS/route/ui.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
	"html/template"

	"github.com/IceWhaleTech/CasaOS/web"
	"github.com/gin-gonic/gin"
)

func WebUIHome(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	index, _ := template.ParseFS(web.Static, "index.html")
	index.Execute(c.Writer, nil)
	return
}
