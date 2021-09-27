package jwt

import (
	"fmt"
	"github.com/IceWhaleTech/CasaOS/model"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT(swagHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = oasis_err2.SUCCESS
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			token = c.Query("token")
		}
		if token == "" {
			code = oasis_err2.INVALID_PARAMS
		}
		if swagHandler == nil {
			claims, err := ParseToken(token)
			if err != nil {
				code = oasis_err2.ERROR_AUTH_TOKEN

			} else if time.Now().Unix() > claims.ExpiresAt {
				code = oasis_err2.ERROR_AUTH_TOKEN
			}
		}

		if code != oasis_err2.SUCCESS {
			c.JSON(http.StatusOK, model.Result{Success: code, Message: oasis_err2.GetMsg(code)})
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
		loger2.NewOLoger().Fatal(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}
