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
	"github.com/IceWhaleTech/CasaOS/codegen/message_bus"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var Cache *cache.Cache

var (
	MyService Repository
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

	MessageBus() *message_bus.ClientWithResponses
	Peer() PeerService
	Other() OtherService
}

func NewService(db *gorm.DB, RuntimePath string) Repository {
	gatewayManagement, err := external.NewManagementService(RuntimePath)
	if err != nil && len(RuntimePath) > 0 {
		panic(err)
	}

	return &store{
		casa:        NewCasaService(),
		connections: NewConnectionsService(db),
		gateway:     gatewayManagement,
		notify:      NewNotifyService(db),
		rely:        NewRelyService(db),
		system:      NewSystemService(),
		health:      NewHealthService(),
		shares:      NewSharesService(db),

		peer:  NewPeerService(db),
		other: NewOtherService(),
	}
}

type store struct {
	peer        PeerService
	db          *gorm.DB
	casa        CasaService
	notify      NotifyServer
	rely        RelyService
	system      SystemService
	shares      SharesService
	connections ConnectionsService
	gateway     external.ManagementService

	health HealthService
	other  OtherService
}

func (c *store) Peer() PeerService {
	return c.peer
}

func (c *store) Other() OtherService {
	return c.other
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

func (c *store) MessageBus() *message_bus.ClientWithResponses {
	client, _ := message_bus.NewClientWithResponses("", func(c *message_bus.Client) error {
		// error will never be returned, as we always want to return a client, even with wrong address,
		// in order to avoid panic.
		//
		// If we don't avoid panic, message bus becomes a hard dependency, which is not what we want.

		messageBusAddress, err := external.GetMessageBusAddress(config.CommonInfo.RuntimePath)
		if err != nil {
			c.Server = "message bus address not found"
			return nil
		}

		c.Server = messageBusAddress
		return nil
	})

	return client
}
