package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/cache"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/sqlite"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/route"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

var sqliteDB *gorm.DB

var configFlag = flag.String("c", "", "config address")
var dbFlag = flag.String("db", "", "db path")
var showUserInfo = flag.Bool("show-user-info", false, "show user info")

func init() {
	flag.Parse()
	config.InitSetup(*configFlag)
	config.UpdateSetup()
	loger2.LogSetup()
	if len(*dbFlag) == 0 {
		*dbFlag = config.AppInfo.ProjectPath + "/db"
	}
	sqliteDB = sqlite.GetDb(*dbFlag)
	//gredis.GetRedisConn(config.RedisInfo),
	service.MyService = service.NewService(sqliteDB, loger2.NewOLoger())
	service.Cache = cache.Init()

	go service.UDPService()

	fmt.Println("token", service.GetToken())
	service.UDPAddressMap = make(map[string]string)
	//go service.SocketConnect()
	service.CancelList = make(map[string]string)
	service.InternalInspection = make(map[string][]string)
	service.NewVersionApp = make(map[string]string)
	route.InitFunction()

	go service.SendIPToServer()
	go service.LoopFriend()
	go service.MyService.App().CheckNewImage()

}

// @title casaOS API
// @version 1.0.0
// @contact.name lauren.pan
// @contact.url https://www.zimaboard.com
// @contact.email lauren.pan@icewhale.org
// @description casaOS v1版本api
// @host 192.168.2.217:8089
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /v1
func main() {
	if *showUserInfo {
		fmt.Println("CasaOS User Info")
		fmt.Println("UserName:" + config.UserInfo.UserName)
		fmt.Println("Password:" + config.UserInfo.PWD)
		return
	}
	//model.Setup()
	//gredis.Setup()
	r := route.InitRouter()
	//service.SyncTask(sqliteDB)
	cron2 := cron.New()
	//every day execution
	err := cron2.AddFunc("0 0/5 * * * *", func() {
		//service.PushIpInfo(*&config.ServerInfo.Token)
		//service.UpdataDDNSList(mysqldb)
		//service.SyncTask(sqliteDB)

		service.SendIPToServer()

		service.LoopFriend()
		service.MyService.App().CheckNewImage()
	})
	if err != nil {
		fmt.Println(err)
	}
	err = cron2.AddFunc("0/1 * * * * *", func() {
		notify := model2.AppNotify{}
		notify.CustomId = ""
		notify.Type = types.NOTIFY_TYPE_HEALTH_CHECK

		go service.MyService.Notify().SendText(notify)

	})
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", config.ServerInfo.HttpPort),
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()

}
