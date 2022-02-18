package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/gorilla/websocket"
)

var WebSocketConn *websocket.Conn

func SocketConnect() {
	Connect()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := WebSocketConn.ReadMessage()
			if err != nil {
				Connect()
			}
			msa := model.MessageModel{}
			json.Unmarshal(message, &msa)
			if msa.Type == "connection" {
				bss, _ := json.Marshal(msa.Data)
				content := model.PersionModel{}
				err := json.Unmarshal(bss, &content)
				fmt.Println(content)
				fmt.Println(err)
			}
		}
	}()

	msg := model.MessageModel{}
	msg.Data = config.ServerInfo.Token
	msg.Type = "refresh"
	b, _ := json.Marshal(msg)
	for {

		select {
		case <-ticker.C:
			err := WebSocketConn.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				Connect()
			}
		case <-done:
			return
		}

	}
}

func Connect() {
	host := strings.Split(config.ServerInfo.Handshake, "://")
	u := url.URL{Scheme: "ws", Host: host[1], Path: "/v1/ws"}

	var err error
	for {
		msg := model.MessageModel{}
		msg.Data = config.ServerInfo.Token
		msg.Type = "join"
		b, _ := json.Marshal(msg)
		if WebSocketConn != nil {
			err = WebSocketConn.WriteMessage(websocket.TextMessage, b)
			if err == nil {
				return
			}
		}

		d, _, e := websocket.DefaultDialer.Dial(u.String(), nil)
		if e == nil {
			WebSocketConn = d
			return
		}
		time.Sleep(time.Second * 5)
	}
}
