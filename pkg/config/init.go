package config

import (
	"fmt"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//系统配置
var SysInfo = &model.SysInfoModel{}

//用户相关
var UserInfo = &model.UserModel{}

//用户相关
var AppInfo = &model.APPModel{}

//redis相关配置
var RedisInfo = &model.RedisModel{}

//zerotier相关
var ZeroTierInfo = &model.ZeroTierModel{}

//server相关
var ServerInfo = &model.ServerModel{}

var SystemConfigInfo = &model.SystemConfig{}

var Cfg *ini.File

//初始化设置，获取系统的部分信息。
func InitSetup(config string) {

	var configDir = USERCONFIGURL
	if len(config) > 0 {
		configDir = config
	}
	var err error
	//读取文件
	Cfg, err = ini.Load(configDir)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mapTo("user", UserInfo)
	mapTo("app", AppInfo)
	mapTo("zerotier", ZeroTierInfo)
	mapTo("redis", RedisInfo)
	mapTo("server", ServerInfo)
	mapTo("system", SystemConfigInfo)
	//	AppInfo.ProjectPath = getCurrentDirectory() //os.Getwd()

}

//映射
func mapTo(section string, v interface{}) {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
