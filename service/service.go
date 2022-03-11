package service

import (
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var Cache *cache.Cache

var MyService Repository

var WebSocketConns []*websocket.Conn

var SocketRun bool

type Repository interface {
	App() AppService
	DDNS() DDNSService
	User() UserService
	Docker() DockerService
	//Redis() RedisService
	ZeroTier() ZeroTierService
	ZiMa() ZiMaService
	Casa() CasaService
	Disk() DiskService
	Notify() NotifyServer
	ShareDirectory() ShareDirService
	Task() TaskService
	Rely() RelyService
	System() SystemService
	Shortcuts() ShortcutsService
	Search() SearchService
	Person() PersonService
	Friend() FriendService
}

func NewService(db *gorm.DB, log loger2.OLog) Repository {

	return &store{
		app:    NewAppService(db, log),
		ddns:   NewDDNSService(db, log),
		user:   NewUserService(),
		docker: NewDockerService(log),
		//redis:      NewRedisService(rp),
		zerotier:       NewZeroTierService(),
		zima:           NewZiMaService(),
		casa:           NewCasaService(),
		disk:           NewDiskService(log, db),
		notify:         NewNotifyService(db),
		shareDirectory: NewShareDirService(db, log),
		task:           NewTaskService(db, log),
		rely:           NewRelyService(db, log),
		system:         NewSystemService(log),
		shortcuts:      NewShortcutsService(db),
		search:         NewSearchService(),
		person:         NewPersonService(db),
		friend:         NewFriendService(db),
	}
}

type store struct {
	db             *gorm.DB
	app            AppService
	ddns           DDNSService
	user           UserService
	docker         DockerService
	zerotier       ZeroTierService
	zima           ZiMaService
	casa           CasaService
	disk           DiskService
	notify         NotifyServer
	shareDirectory ShareDirService
	task           TaskService
	rely           RelyService
	system         SystemService
	shortcuts      ShortcutsService
	search         SearchService
	person         PersonService
	friend         FriendService
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

func (c *store) DDNS() DDNSService {
	return c.ddns
}

func (c *store) User() UserService {
	return c.user
}

func (c *store) Docker() DockerService {
	return c.docker
}

func (c *store) ZeroTier() ZeroTierService {
	return c.zerotier
}
func (c *store) ZiMa() ZiMaService {
	return c.zima
}
func (c *store) Casa() CasaService {
	return c.casa
}

func (c *store) Disk() DiskService {
	return c.disk
}
func (c *store) ShareDirectory() ShareDirService {
	return c.shareDirectory
}
func (c *store) Task() TaskService {
	return c.task
}
func (c *store) Search() SearchService {
	return c.search
}
