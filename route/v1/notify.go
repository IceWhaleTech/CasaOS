package v1

import (
	json2 "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"oasis/model"
	"oasis/pkg/utils/oasis_err"
	"oasis/service"
	"oasis/types"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary websocket 接口,连接成功后发送一个"notify"字符串
// @Produce  application/json
// @Accept application/json
// @Tags notify
// @Security ApiKeyAuth
// @Param token path string true "token"
// @Success 200 {string} string "ok"
// @Router /notify/ws [get]
func NotifyWS(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) != "notify" {
			return
		}
		for {
			list := service.MyService.Notify().GetList()
			json, _ := json2.Marshal(list)
			err = ws.WriteMessage(mt, json)
			if err != nil {
				break
			}
			time.Sleep(time.Second * 2)
		}
	}
}

// @Summary 标记notify已读
// @Produce  application/json
// @Accept application/json
// @Tags notify
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /notify/read/{id} [put]
func PutNotifyRead(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	service.MyService.Notify().MarkRead(id, types.NOTIFY_READ)
}
