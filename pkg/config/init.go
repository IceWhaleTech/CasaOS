/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-05 11:58:02
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
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS/common"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/go-ini/ini"
)

var (
	SysInfo = &model.SysInfoModel{}
	AppInfo = &model.APPModel{
		DBPath:       constants.DefaultDataPath,
		LogPath:      constants.DefaultLogPath,
		LogSaveName:  common.SERVICENAME,
		LogFileExt:   "log",
		ShellPath:    "/usr/share/casaos/shell",
		UserDataPath: filepath.Join(constants.DefaultDataPath, "conf"),
	}
	CommonInfo = &model.CommonModel{
		RuntimePath: constants.DefaultRuntimePath,
	}
	ServerInfo       = &model.ServerModel{}
	SystemConfigInfo = &model.SystemConfig{}
	FileSettingInfo  = &model.FileSetting{}

	Cfg            *ini.File
	ConfigFilePath string
)

// 初始化设置，获取系统的部分信息。
func InitSetup(config string, sample string) {
	ConfigFilePath = CasaOSConfigFilePath
	if len(config) > 0 {
		ConfigFilePath = config
	}

	// create default config file if not exist
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		fmt.Println("config file not exist, create it")
		// create config file
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write default config
		_, err = file.WriteString(sample)
		if err != nil {
			panic(err)
		}
	}

	var err error

	// 读取文件
	Cfg, err = ini.Load(ConfigFilePath)
	if err != nil {
		panic(err)
	}

	mapTo("app", AppInfo)
	mapTo("server", ServerInfo)
	mapTo("system", SystemConfigInfo)
	mapTo("file", FileSettingInfo)
	mapTo("common", CommonInfo)
}

// 映射
func mapTo(section string, v interface{}) {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
