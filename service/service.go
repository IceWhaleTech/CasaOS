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
	DDNS() DDNSService
	User() UserService
	Docker() DockerService
	//Redis() RedisService
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
	Download() DownloadService
	DownRecord() DownRecordService
}

func NewService(db *gorm.DB) Repository {

	return &store{
		app:    NewAppService(db),
		user:   NewUserService(),
		docker: NewDockerService(),
		//redis:      NewRedisService(rp),
		zima:           NewZiMaService(),
		casa:           NewCasaService(),
		disk:           NewDiskService(db),
		notify:         NewNotifyService(db),
		shareDirectory: NewShareDirService(db),
		task:           NewTaskService(db),
		rely:           NewRelyService(db),
		system:         NewSystemService(),
		shortcuts:      NewShortcutsService(db),
		search:         NewSearchService(),
		person:         NewPersonService(db),
		friend:         NewFriendService(db),
		download:       NewDownloadService(db),
		downrecord:     NewDownRecordService(db),
	}
}

type store struct {
	db             *gorm.DB
	app            AppService
	ddns           DDNSService
	user           UserService
	docker         DockerService
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
	download       DownloadService
	downrecord     DownRecordService
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

func (c *store) DDNS() DDNSService {
	return c.ddns
}

func (c *store) User() UserService {
	return c.user
}

func (c *store) Docker() DockerService {
	return c.docker
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
