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
	gateway "github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var Cache *cache.Cache

var MyService Repository

var WebSocketConns []*websocket.Conn
var NewVersionApp map[string]string
var SocketRun bool

type Repository interface {
	App() AppService
	//User() UserService
	Docker() DockerService
	Casa() CasaService
	Notify() NotifyServer
	Rely() RelyService
	System() SystemService
	Shares() SharesService
	Connections() ConnectionsService
	Gateway() gateway.ManagementService
}

func NewService(db *gorm.DB, RuntimePath string) Repository {

	gatewayManagement, err := gateway.NewManagementService(RuntimePath)
	if err != nil && len(RuntimePath) > 0 {
		panic(err)
	}

	return &store{
		gateway:     gatewayManagement,
		app:         NewAppService(db),
		docker:      NewDockerService(),
		casa:        NewCasaService(),
		notify:      NewNotifyService(db),
		rely:        NewRelyService(db),
		system:      NewSystemService(),
		shares:      NewSharesService(db),
		connections: NewConnectionsService(db),
	}
}

type store struct {
	db          *gorm.DB
	app         AppService
	docker      DockerService
	casa        CasaService
	notify      NotifyServer
	rely        RelyService
	system      SystemService
	shares      SharesService
	connections ConnectionsService
	gateway     gateway.ManagementService
}

func (c *store) Gateway() gateway.ManagementService {
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

func (c *store) App() AppService {
	return c.app
}

func (c *store) Docker() DockerService {
	return c.docker
}

func (c *store) Casa() CasaService {
	return c.casa
}
