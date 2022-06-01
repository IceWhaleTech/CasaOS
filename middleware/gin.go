/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-10-08 10:29:08
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-25 19:17:45
 * @FilePath: /CasaOS/middleware/gin.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,Language")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Set("content-type", "application/json")
		//}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		c.Next()
	}
}
