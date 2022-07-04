/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-17 14:01:25
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-04 16:26:22
 * @FilePath: /CasaOS/pkg/utils/jwt/jwt_helper.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package jwt

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/gin-gonic/gin"
)

// func JWT() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var code int
// 		code = common_err.SUCCESS
// 		token := c.GetHeader("Authorization")
// 		if len(token) == 0 {
// 			token = c.Query("token")
// 		}
// 		if token == "" {
// 			code = common_err.INVALID_PARAMS
// 		}

// 		claims, err := ParseToken(token)
// 		//_, err := ParseToken(token)
// 		if err != nil {
// 			code = common_err.ERROR_AUTH_TOKEN
// 		} else if !claims.VerifyExpiresAt(time.Now(), true) || !claims.VerifyIssuer("casaos", true) {
// 			code = common_err.ERROR_AUTH_TOKEN
// 		}
// 		if code != common_err.SUCCESS {
// 			c.JSON(code, model.Result{Success: code, Message: common_err.GetMsg(code)})
// 			c.Abort()
// 			return
// 		}
// 		c.Request.Header.Add("user_id", strconv.Itoa(claims.Id))
// 		c.Next()
// 	}
// }

// //get AccessToken
// func GetAccessToken(username, pwd string, id int) string {
// 	token, err := GenerateToken(username, pwd, id, "casaos", 3*time.Hour*time.Duration(1))
// 	if err == nil {
// 		return token
// 	} else {
// 		loger2.Error(fmt.Sprintf("Get Token Fail: %V", err))
// 		return ""
// 	}
// }

// func GetRefreshToken(username, pwd string, id int) string {
// 	token, err := GenerateToken(username, pwd, id, "fresh", 7*24*time.Hour*time.Duration(1))
// 	if err == nil {
// 		return token
// 	} else {
// 		loger2.Error(fmt.Sprintf("Get Token Fail: %V", err))
// 		return ""
// 	}
// }

//*************** soon to be removed *****************

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

		//claims, err := ParseToken(token)
		_, err := ParseToken(token)
		if err != nil {
			code = common_err.ERROR_AUTH_TOKEN

		}
		//else if time.Now().Unix() > claims.ExpiresAt {
		//	code = oasis_err2.ERROR_AUTH_TOKEN
		//}

		if code != common_err.SUCCESS {
			c.JSON(http.StatusOK, model.Result{Success: code, Message: common_err.GetMsg(code)})
			c.Abort()
			return
		}
		c.Next()
	}
}

//获取token
func GetToken(username, pwd string) string {
	token, err := GenerateToken(username, pwd)
	if err == nil {
		return token
	} else {
		//loger2.NewOLoger().Fatal(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}
