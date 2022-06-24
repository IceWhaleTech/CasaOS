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
	Shortcuts() ShortcutsService
	Person() PersonService
	Friend() FriendService
	Download() DownloadService
	DownRecord() DownRecordService
}

func NewService(db *gorm.DB) Repository {

	return &store{
		app:        NewAppService(db),
		user:       NewUserService(db),
		docker:     NewDockerService(),
		casa:       NewCasaService(),
		disk:       NewDiskService(db),
		notify:     NewNotifyService(db),
		rely:       NewRelyService(db),
		system:     NewSystemService(),
		shortcuts:  NewShortcutsService(db),
		person:     NewPersonService(db),
		friend:     NewFriendService(db),
		download:   NewDownloadService(db),
		downrecord: NewDownRecordService(db),
	}
}

type store struct {
	db         *gorm.DB
	app        AppService
	user       UserService
	docker     DockerService
	casa       CasaService
	disk       DiskService
	notify     NotifyServer
	rely       RelyService
	system     SystemService
	shortcuts  ShortcutsService
	person     PersonService
	friend     FriendService
	download   DownloadService
	downrecord DownRecordService
}

func (c *store) DownRecord() DownRecordService {
	return c.downrecord
}

func (c *store) Download() DownloadService {
	return c.download
}
func (c *store) Friend() FriendService {
	return c.friend
}
func (c *store) Rely() RelyService {
	return c.rely
}
func (c *store) Shortcuts() ShortcutsService {
	return c.shortcuts
}
func (c *store) Person() PersonService {
	return c.person
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
