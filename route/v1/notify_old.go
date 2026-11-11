package v1

import (
	"net/http"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		return strings.HasPrefix(origin, "http://localhost") ||
			strings.HasPrefix(origin, "http://127.0.0.1") ||
			strings.HasPrefix(origin, "https://localhost") ||
			strings.HasPrefix(origin, "https://127.0.0.1") ||
			strings.HasPrefix(origin, "http://192.168.") ||
			strings.HasPrefix(origin, "http://10.") ||
			strings.HasPrefix(origin, "http://172.")
	},
}

// @Summary websocket endpoint, sends "notify" string on successful connection
// @Produce  application/json
// @Accept application/json
// @Tags notify
// @Security ApiKeyAuth
// @Param token path string true "token"
// @Success 200 {string} string "ok"
// @Router /notify/ws [get]
func NotifyWS(ctx echo.Context) error {
	ws, err := upGrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		logger.Error("failed to upgrade websocket connection", zap.Error(err))
		return nil
	}
	defer ws.Close()
	service.WebSocketConns = append(service.WebSocketConns, ws)

	if !service.SocketRun {
		service.SocketRun = true
		service.SendMeg()
	}
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("websocket read error", zap.Error(err))
			}
			break
		}
	}
	return nil
}

// @Summary Mark notify as read
// @Produce  application/json
// @Accept application/json
// @Tags notify
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /notify/read/{id} [put]
func PutNotifyRead(ctx echo.Context) error {
	id := ctx.Param("id")
	service.MyService.Notify().MarkRead(id, types.NOTIFY_READ)
	return nil
}
