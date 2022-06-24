/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-21 16:01:26
 * @FilePath: /CasaOS/pkg/config/init.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/go-ini/ini"
)

//系统配置
var SysInfo = &model.SysInfoModel{}

//用户相关
var UserInfo = &model.UserModel{}

//用户相关
var AppInfo = &model.APPModel{}

//var RedisInfo = &model.RedisModel{}

//server相关
var ServerInfo = &model.ServerModel{}

var SystemConfigInfo = &model.SystemConfig{}

var CasaOSGlobalVariables = &model.CasaOSGlobalVariables{}

var FileSettingInfo = &model.FileSetting{}

var Cfg *ini.File

//初始化设置，获取系统的部分信息。
func InitSetup(config string) {

	var configDir = USERCONFIGURL
	if len(config) > 0 {
		configDir = config
	}
	if runtime.GOOS == "darwin" {
		configDir = "./conf/conf.conf"
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
	//mapTo("redis", RedisInfo)
	mapTo("server", ServerInfo)
	mapTo("system", SystemConfigInfo)
	mapTo("file", FileSettingInfo)
	SystemConfigInfo.ConfigPath = configDir
	if len(AppInfo.DBPath) == 0 {
		AppInfo.DBPath = "/var/lib/casaos"
		Cfg.SaveTo(configDir)
	}
	if len(AppInfo.LogPath) == 0 {
		AppInfo.LogPath = "/var/log/casaos/"
		Cfg.SaveTo(configDir)
	}
	if len(AppInfo.ShellPath) == 0 {
		AppInfo.ShellPath = "/usr/share/casaos/shell"
		Cfg.SaveTo(configDir)
	}
	if len(AppInfo.UserDataPath) == 0 {
		AppInfo.UserDataPath = "/var/lib/casaos/conf"
		Cfg.SaveTo(configDir)
	}
	if len(AppInfo.TempPath) == 0 {
		AppInfo.TempPath = "/var/lib/casaos/temp"
		Cfg.SaveTo(configDir)
	}
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
