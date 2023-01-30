/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-12 09:48:56
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-02 22:10:05
 * @FilePath: /CasaOS/service/service.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/external"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Cache *cache.Cache

var (
	MyService    Repository
	SocketServer *socketio.Server
)

var (
	WebSocketConns []*websocket.Conn
	SocketRun      bool
)

type Repository interface {
	Casa() CasaService
	Connections() ConnectionsService
	Gateway() external.ManagementService
	Health() HealthService
	Notify() NotifyServer
	Rely() RelyService
	Shares() SharesService
	System() SystemService
}

func NewService(db *gorm.DB, RuntimePath string, socket *socketio.Server) Repository {
	if socket == nil {
		logger.Error("socket is nil", zap.Any("error", "socket is nil"))
	}
	SocketServer = socket
	gatewayManagement, err := external.NewManagementService(RuntimePath)
	if err != nil && len(RuntimePath) > 0 {
		panic(err)
	}

	return &store{
		casa:        NewCasaService(),
		connections: NewConnectionsService(db),
		gateway:     gatewayManagement,
		health:      NewHealthService(),
		notify:      NewNotifyService(db),
		rely:        NewRelyService(db),
		shares:      NewSharesService(db),
		system:      NewSystemService(),
	}
}

type store struct {
	db          *gorm.DB
	casa        CasaService
	connections ConnectionsService
	gateway     external.ManagementService
	health      HealthService
	notify      NotifyServer
	rely        RelyService
	shares      SharesService
	system      SystemService
}

func (c *store) Gateway() external.ManagementService {
	return c.gateway
}

func (s *store) Connections() ConnectionsService {
	return s.connections
}

func (s *store) Shares() SharesService {
	return s.shares
}

func (c *store) Rely() RelyService {
	return c.rely
}

func (c *store) System() SystemService {
	return c.system
}

func (c *store) Notify() NotifyServer {
	return c.notify
}

func (c *store) Casa() CasaService {
	return c.casa
}

func (c *store) Health() HealthService {
	return c.health
}
