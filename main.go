package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/pkg/cache"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/sqlite"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/route"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/IceWhaleTech/CasaOS/types"
	"go.uber.org/zap"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

const LOCALHOST = "127.0.0.1"

var sqliteDB *gorm.DB

var (
	configFlag  = flag.String("c", "", "config address")
	dbFlag      = flag.String("db", "", "db path")
	versionFlag = flag.Bool("v", false, "version")
)

func init() {
	flag.Parse()
	if *versionFlag {
		fmt.Println("v" + types.CURRENTVERSION)
		return
	}
	config.InitSetup(*configFlag)
	config.UpdateSetup()

	loger.LogInit()
	if len(*dbFlag) == 0 {
		*dbFlag = config.AppInfo.DBPath + "/db"
	}

	sqliteDB = sqlite.GetDb(*dbFlag)
	// gredis.GetRedisConn(config.RedisInfo),

	service.MyService = service.NewService(sqliteDB, config.CommonInfo.RuntimePath)

	service.Cache = cache.Init()

	service.GetToken()

	service.NewVersionApp = make(map[string]string)
	route.InitFunction()

	// go service.LoopFriend()
	// go service.MyService.App().CheckNewImage()
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
	service.NotifyMsg = make(chan notify.Message, 10)
	if *versionFlag {
		return
	}
	go route.SocketInit(service.NotifyMsg)
	// model.Setup()
	// gredis.Setup()

	r := route.InitRouter()
	// service.SyncTask(sqliteDB)
	cron2 := cron.New()
	// every day execution

	err := cron2.AddFunc("0/5 * * * * *", func() {
		if service.ClientCount > 0 {
			// route.SendNetINfoBySocket()
			// route.SendCPUBySocket()
			// route.SendMemBySocket()
			// route.SendDiskBySocket()
			// route.SendUSBBySocket()
			route.SendAllHardwareStatusBySocket()
		}
	})
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()

	defer cron2.Stop()

	listener, err := net.Listen("tcp", net.JoinHostPort(LOCALHOST, "0"))
	if err != nil {
		panic(err)
	}
	routers := []string{"sys", "apps", "container", "app-categories", "port", "file", "folder", "batch", "image", "samba", "notify"}
	for _, v := range routers {
		err = service.MyService.Gateway().CreateRoute(&common.Route{
			Path:   "/v1/" + v,
			Target: "http://" + listener.Addr().String(),
		})

		if err != nil {
			fmt.Println("err", err)
			panic(err)
		}
	}
	go func() {
		time.Sleep(time.Second * 2)
		// v0.3.6
		if config.ServerInfo.HttpPort != "" {
			changePort := common.ChangePortRequest{}
			changePort.Port = config.ServerInfo.HttpPort
			err := service.MyService.Gateway().ChangePort(&changePort)
			if err == nil {
				config.Cfg.Section("server").Key("HttpPort").SetValue("")
				config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
			}
		}
	}()

	// s := &http.Server{
	// 	Addr:           listener.Addr().String(), //fmt.Sprintf(":%v", config.ServerInfo.HttpPort),
	// 	Handler:        r,
	// 	ReadTimeout:    60 * time.Second,
	// 	WriteTimeout:   60 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.ListenAndServe()
	urlFilePath := filepath.Join(config.CommonInfo.RuntimePath, "casaos.url")
	err = file.CreateFileAndWriteContent(urlFilePath, "http://"+listener.Addr().String())
	if err != nil {
		loger.Error("Management service is listening...",
			zap.Any("address", listener.Addr().String()),
			zap.Any("filepath", urlFilePath),
		)
	}

	err = http.Serve(listener, r)
	if err != nil {
		panic(err)
	}
}
