package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
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
				//开始尝试udp链接
				go UDPConnect(content.Ips)
			} else if msa.Type == types.PERSONADDFRIEND {
				// new add friend
				uuid := uuid.NewV4().String()
				mi := model2.FriendModel{}
				mi.Avatar = config.UserInfo.Avatar
				mi.Profile = config.UserInfo.Description
				mi.Name = config.UserInfo.UserName
				m := model.MessageModel{}
				m.Data = mi
				m.From = config.ServerInfo.Token
				m.To = msa.From
				m.Type = types.PERSONADDFRIEND
				m.UUId = uuid
				result, err := Dial("192.168.2.225:9902", m)
				friend := model2.FriendModel{}
				if err != nil && !reflect.DeepEqual(result, friend) {
					dataModelByte, _ := json.Marshal(result.Data)
					err := json.Unmarshal(dataModelByte, &friend)
					if err != nil {
						fmt.Println(err)
					}
				}
				if len(friend.Token) == 0 {
					friend.Token = m.From
				}
				MyService.Friend().AddFriend(friend)
			}
		}
	}()

	msg := model.MessageModel{}
	msg.Data = config.ServerInfo.Token
	msg.Type = "refresh"
	msg.From = config.ServerInfo.Token
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
	for {
		d, _, e := websocket.DefaultDialer.Dial(u.String(), nil)
		if e == nil {
			WebSocketConn = d
			return
		}
		time.Sleep(time.Second * 5)
	}
}
