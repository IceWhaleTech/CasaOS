package v1

import (
	"fmt"
	"net/http"

	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
func NotifyWS(ctx echo.Context) error {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		return nil
	}
	defer ws.Close()
	service.WebSocketConns = append(service.WebSocketConns, ws)

	if !service.SocketRun {
		service.SocketRun = true
		service.SendMeg()
	}
	for {
		mt, message, err := ws.ReadMessage()
		fmt.Println(mt, message, err)
	}
}

// @Summary 标记notify已读
// @Produce  application/json
// @Accept application/json
// @Tags notify
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /notify/read/{id} [put]
func PutNotifyRead(ctx echo.Context) error {
	id := ctx.Param("id")
	// if len(id) == 0 {
	// 	return ctx.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
	// 	return
	// }
	fmt.Println(id)
	service.MyService.Notify().MarkRead(id, types.NOTIFY_READ)
	return nil
}
