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

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS-Common/utils/port"
	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/service"
	f "github.com/ambelovsky/gosf"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
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
		"port": socketPort,
	})
}

var SocketServer *socketio.Server

func SocketIo() *socketio.Server {
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("connected", zap.Any("id", s.ID()))
		s.Join("bcast")
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		logger.Info("notice", zap.Any("msg", msg))
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("meet error", zap.Any("error", e))
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Info("closed", zap.Any("reason", reason))
	})

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error("error when trying to  listen socketio ", zap.Any("error", err))
		}
	}()
	SocketServer = server
	return server
}

func SendDataBySocketIo(path string, data interface{}) {
	if SocketServer != nil {
		SocketServer.BroadcastToRoom("/", "bcast", path, data)
	}
}
