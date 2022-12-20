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
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/service"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
)

func SocketIo() *socketio.Server {
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("connected", zap.Any("id", s.ID()))
		s.Join("public")
		service.ClientCount += 1
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		logger.Info("notice", zap.Any("msg", msg))
		s.Emit("reply", "have "+msg)
	})

	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
	// 	last := s.Context().(string)
	// 	s.Emit("bye", last)
	// 	s.Close()
	// 	return last
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("meet error", zap.Any("error", e))
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		service.ClientCount -= 1
		logger.Info("closed", zap.Any("reason", reason))
	})

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error("error when trying to  listen socketio ", zap.Any("error", err))
		}
	}()
	return server
}
