package jwt

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/gin-gonic/gin"
)

func JWT(swagHandler gin.HandlerFunc) gin.HandlerFunc {
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
		if swagHandler == nil {
			claims, err := ParseToken(token)
			//_, err := ParseToken(token)
			if err != nil {
				code = common_err.ERROR_AUTH_TOKEN
			} else if claims.VerifyExpiresAt(time.Now(), true) || claims.VerifyIssuer("casaos", true) {
				code = common_err.ERROR_AUTH_TOKEN
			}
			c.Header("user_id", strconv.Itoa(claims.Id))
		}

		if code != common_err.SUCCESS {
			c.JSON(http.StatusOK, model.Result{Success: code, Message: common_err.GetMsg(code)})
			c.Abort()
			return
		}

		c.Next()
	}
}

//get AccessToken
func GetAccessToken(username, pwd string, id int) string {
	token, err := GenerateToken(username, pwd, id, "casaos", 3*time.Hour*time.Duration(1))
	if err == nil {
		return token
	} else {
		loger2.NewOLoger().Fatal(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}

func GetRefreshToken(username, pwd string, id int) string {
	token, err := GenerateToken(username, pwd, id, "fresh", 7*24*time.Hour*time.Duration(1))
	if err == nil {
		return token
	} else {
		loger2.NewOLoger().Fatal(fmt.Sprintf("Get Token Fail: %V", err))
		return ""
	}
}
