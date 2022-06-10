/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-23 17:18:56
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-09 21:48:10
 * @FilePath: /CasaOS/route/socket.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/service"
	f "github.com/ambelovsky/gosf"
)

func SocketInit(msg chan notify.Message) {

	// set socket port
	socketPort := 0
	if len(config.ServerInfo.SocketPort) == 0 {
		socketPort, _ = port.GetAvailablePort("tcp")
		config.ServerInfo.SocketPort = strconv.Itoa(socketPort)
		config.Cfg.Section("server").Key("SocketPort").SetValue(strconv.Itoa(socketPort))
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	} else {
		socketPort, _ = strconv.Atoi(config.ServerInfo.SocketPort)
		if !port.IsPortAvailable(socketPort, "tcp") {
			socketPort, _ := port.GetAvailablePort("tcp")
			config.ServerInfo.SocketPort = strconv.Itoa(socketPort)
			config.Cfg.Section("server").Key("SocketPort").SetValue(strconv.Itoa(socketPort))
			config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
		}
	}

	f.OnConnect(func(c *f.Client, request *f.Request) {
		service.ClientCount += 1
	})
	f.OnDisconnect(func(c *f.Client, request *f.Request) {
		service.ClientCount -= 1
	})
	go func(msg chan notify.Message) {
		for v := range msg {
			f.Broadcast("", v.Path, &v.Msg)
			time.Sleep(time.Millisecond * 100)
		}

	}(msg)

	f.Startup(map[string]interface{}{
		"port": socketPort})

}
