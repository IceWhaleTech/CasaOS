package main

import (
	"flag"
	"fmt"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/sqlite"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/route"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var sqliteDB *gorm.DB

var swagHandler gin.HandlerFunc
var configFlag = flag.String("c", "", "config address")

func init() {
	flag.Parse()
	config.InitSetup(*configFlag)
	loger2.LogSetup()
	sqliteDB = sqlite.GetDb(config.AppInfo.ProjectPath)
	//gredis.GetRedisConn(config.RedisInfo),
	service.MyService = service.NewService(sqliteDB, loger2.NewOLoger())
}

// @title Oasis API
// @version 1.0.0
// @contact.name lauren.pan
// @contact.url https://www.zimaboard.com
// @contact.email lauren.pan@icewhale.org
// @description Oasis v1版本api
// @host 192.168.2.114:8089
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /v1
func main() {
	//model.Setup()
	//gredis.Setup()
	r := route.InitRouter(swagHandler)
	service.SyncTask(sqliteDB)
	cron2 := cron.New() //创建一个cron实例
	//执行定时任务（每5秒执行一次）
	err := cron2.AddFunc("0 0 0 1/1 * *", func() {
		//service.UpdataDDNSList(mysqldb)
		service.SyncTask(sqliteDB)
	})
	if err != nil {
		fmt.Println(err)
	}

	//启动/关闭
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
