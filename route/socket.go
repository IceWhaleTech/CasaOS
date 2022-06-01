/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-23 17:18:56
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-30 17:06:08
 * @FilePath: /CasaOS/route/socket.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
	"time"

	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/service"
	f "github.com/ambelovsky/gosf"
)

func ScoketInit(port int, msg chan notify.Message) {
	f.OnConnect(func(c *f.Client, request *f.Request) {
		service.ClientCount += 1
	})
	f.OnDisconnect(func(c *f.Client, request *f.Request) {
		service.ClientCount -= 1
	})
	go func(msg chan notify.Message) {
		for v := range msg {
			f.Broadcast("", v.Path, &v.Msg)
			time.Sleep(time.Millisecond * 300)
		}

	}(msg)

	f.Startup(map[string]interface{}{
		"port": port})

}
