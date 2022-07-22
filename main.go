package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS/model/notify"
	"github.com/IceWhaleTech/CasaOS/pkg/cache"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/sqlite"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/random"
	"github.com/IceWhaleTech/CasaOS/route"
	"github.com/IceWhaleTech/CasaOS/service"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

var sqliteDB *gorm.DB

var configFlag = flag.String("c", "", "config address")
var dbFlag = flag.String("db", "", "db path")
var resetUser = flag.Bool("ru", false, "reset user")
var user = flag.String("user", "", "user name")

func init() {
	flag.Parse()
	config.InitSetup(*configFlag)
	config.UpdateSetup()
	loger.LogInit()
	if len(*dbFlag) == 0 {
		*dbFlag = config.AppInfo.DBPath + "/db"
	}
	sqliteDB = sqlite.GetDb(*dbFlag)
	//gredis.GetRedisConn(config.RedisInfo),
	service.MyService = service.NewService(sqliteDB)
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
	if *resetUser {
		if user == nil || len(*user) == 0 {
			fmt.Println("user is empty")
			return
		}
		userData := service.MyService.User().GetUserAllInfoByName(*user)
		if userData.Id == 0 {
			fmt.Println("user not exist")
			return
		}
		password := random.RandomString(6, false)
		userData.Password = encryption.GetMD5ByStr(password)
		service.MyService.User().UpdateUserPassword(userData)
		fmt.Println("User reset successful")
		fmt.Println("UserName:" + userData.Username)
		fmt.Println("Password:" + password)
		return
	}
	go route.SocketInit(service.NotifyMsg)

	//model.Setup()
	//gredis.Setup()
	r := route.InitRouter()
	//service.SyncTask(sqliteDB)
	cron2 := cron.New()
	//every day execution

	err := cron2.AddFunc("0/5 * * * * *", func() {
		if service.ClientCount > 0 {
			//route.SendNetINfoBySocket()
			//route.SendCPUBySocket()
			//route.SendMemBySocket()
			// route.SendDiskBySocket()
			//route.SendUSBBySocket()
			route.SendAllHardwareStatusBySocket()
		}
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

	// if err := r.Run(fmt.Sprintf(":%v", config.ServerInfo.HttpPort)); err != nil {
	// 	fmt.Println("failed run app: ", err)
	// }
}
