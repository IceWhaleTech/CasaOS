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

// InitSetup initializes settings and retrieves part of the system information.
func InitSetup(config, sample string) error {
	ConfigFilePath = CasaOSConfigFilePath
	if len(config) > 0 {
		ConfigFilePath = config
	}

	// Atomically create config file if it doesn't exist (O_EXCL prevents race condition)
	file, err := os.OpenFile(ConfigFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err == nil {
		// File was created, write default config
		defer file.Close()
		if _, err = file.WriteString(sample); err != nil {
			return fmt.Errorf("failed to write default config: %w", err)
		}
	} else if !os.IsExist(err) {
		// Error other than "file exists"
		return fmt.Errorf("failed to create config file: %w", err)
	}

	Cfg, err = ini.Load(ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	if err := mapTo("app", AppInfo); err != nil {
		return err
	}
	if err := mapTo("server", ServerInfo); err != nil {
		return err
	}
	if err := mapTo("system", SystemConfigInfo); err != nil {
		return err
	}
	if err := mapTo("file", FileSettingInfo); err != nil {
		return err
	}
	if err := mapTo("common", CommonInfo); err != nil {
		return err
	}

	return nil
}

func mapTo(section string, v interface{}) error {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		return fmt.Errorf("failed to map config section %s: %w", section, err)
	}
	return nil
}
