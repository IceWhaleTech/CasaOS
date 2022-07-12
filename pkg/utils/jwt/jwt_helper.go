/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-17 14:01:25
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-11 18:21:40
 * @FilePath: /CasaOS/pkg/utils/jwt/jwt_helper.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = common_err.SUCCESS
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			token = c.Query("token")
		}
		if token == "" {
			code = common_err.INVALID_PARAMS
		}

		claims, err := ParseToken(token)
		//_, err := ParseToken(token)
		if err != nil {
			code = common_err.ERROR_AUTH_TOKEN
		} else if !claims.VerifyExpiresAt(time.Now(), true) || !claims.VerifyIssuer("casaos", true) {
			code = common_err.ERROR_AUTH_TOKEN
		}
		if code != common_err.SUCCESS {
			c.JSON(code, model.Result{Success: code, Message: common_err.GetMsg(code)})
			c.Abort()
			return
		}
		c.Request.Header.Add("user_id", strconv.Itoa(claims.Id))
		c.Next()
	}
}

//get AccessToken
func GetAccessToken(username, pwd string, id int) string {
	token, err := GenerateToken(username, pwd, id, "casaos", 3*time.Hour*time.Duration(1))
	if err == nil {
		return token
	} else {
		loger2.Error(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}

func GetRefreshToken(username, pwd string, id int) string {
	token, err := GenerateToken(username, pwd, id, "refresh", 7*24*time.Hour*time.Duration(1))
	if err == nil {
		return token
	} else {
		loger2.Error(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}
