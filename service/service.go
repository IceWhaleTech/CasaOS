/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-12 09:48:56
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-15 10:58:54
 * @FilePath: /CasaOS/service/service.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
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
	User() UserService
	Docker() DockerService
	Casa() CasaService
	Disk() DiskService
	Notify() NotifyServer
	Rely() RelyService
	System() SystemService
}

func NewService(db *gorm.DB) Repository {
	return &store{
		app:    NewAppService(db),
		user:   NewUserService(db),
		docker: NewDockerService(),
		casa:   NewCasaService(),
		disk:   NewDiskService(db),
		notify: NewNotifyService(db),
		rely:   NewRelyService(db),
		system: NewSystemService(),
	}
}

type store struct {
	db     *gorm.DB
	app    AppService
	user   UserService
	docker DockerService
	casa   CasaService
	disk   DiskService
	notify NotifyServer
	rely   RelyService
	system SystemService
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

func (c *store) User() UserService {
	return c.user
}

func (c *store) Docker() DockerService {
	return c.docker
}

func (c *store) Casa() CasaService {
	return c.casa
}

func (c *store) Disk() DiskService {
	return c.disk
}
