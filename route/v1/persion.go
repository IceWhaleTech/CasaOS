package v1

import (
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func PersonTest(c *gin.Context) {

	//service.MyService.Person().GetPersionInfo("fb2333a1-72b2-4cb4-9e31-61ccaffa55b9")

	m := model.ConnectState{}
	m.CreatedAt = time.Now()
	m.From = config.ServerInfo.Token
	m.To = "fb2333a1-72b2-4cb4-9e31-61ccaffa55b9"
	m.Type = ""
	m.UUId = uuid.NewV4().String()

	//service.MyService.Person().Handshake(m)
	err := service.WebSocketConn.WriteMessage(websocket.TextMessage, []byte("test1111"))
	if err == nil {
		return
	}
}
