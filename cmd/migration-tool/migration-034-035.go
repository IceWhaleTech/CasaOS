/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-08-24 17:36:00
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-05 11:24:27
 * @FilePath: /CasaOS/cmd/migration-tool/migration-034-035.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"
)

type migrationTool struct{}

func (u *migrationTool) IsMigrationNeeded() (bool, error) {

	majorVersion, minorVersion, patchVersion, err := version.DetectLegacyVersion()
	if err != nil {
		if err == version.ErrLegacyVersionNotFound {
			return false, nil
		}

		return false, err
	}

	if majorVersion > 0 {
		return false, nil
	}

	if minorVersion > 3 {
		return false, nil
	}

	if minorVersion == 3 && patchVersion > 5 {
		return false, nil
	}

	_logger.Info("Migration is needed for a CasaOS version 0.3.5 and older...")
	return true, nil

}

func (u *migrationTool) PreMigrate() error {

	return nil
}

func (u *migrationTool) Migrate() error {

	if service.MyService.System().GetSysInfo().KernelArch == "aarch64" && config.ServerInfo.USBAutoMount != "True" && strings.Contains(service.MyService.System().GetDeviceTree(), "Raspberry Pi") {
		service.MyService.System().UpdateUSBAutoMount("False")
		service.MyService.System().ExecUSBAutoMountShell("False")
	}
	newAPIUrl := "https://api.casaos.io/casaos-api"
	if config.ServerInfo.ServerApi == "https://api.casaos.zimaboard.com" {
		config.ServerInfo.ServerApi = newAPIUrl
		config.Cfg.Section("server").Key("ServerApi").SetValue(newAPIUrl)
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}
	command.OnlyExec("curl -fsSL https://raw.githubusercontent.com/IceWhaleTech/get/main/assist.sh | bash")
	if !file.CheckNotExist("/casaOS") {
		command.OnlyExec("source /casaOS/server/shell/update.sh ;")
		command.OnlyExec("source " + config.AppInfo.ShellPath + "/delete-old-service.sh ;")
	}

	service.MyService.App().ImportApplications(true)

	src := "/casaOS/server/conf/conf.ini"
	if file.Exists(src) {
		dst := "/etc/casaos/casaos.conf"
		source, err := os.Open(src)
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}

	}

	if file.Exists("/casaOS/server/db") {
		var fds []os.FileInfo
		var err error
		to := "/var/lib/casaos/db"
		file.IsNotExistMkDir(to)
		from := "/casaOS/server/db"
		if fds, err = ioutil.ReadDir(from); err != nil {
			return err
		}

		for _, fd := range fds {
			srcfp := path.Join(from, fd.Name())
			dstfp := path.Join(to, fd.Name())
			source, err := os.Open(srcfp)
			if err != nil {
				return err
			}
			defer source.Close()

			destination, err := os.Create(dstfp)
			if err != nil {
				return err
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				return err
			}
		}

	}

	if file.Exists("/casaOS/server/conf") {
		var fds []os.FileInfo
		var err error
		to := "/var/lib/casaos/conf"
		file.IsNotExistMkDir(to)
		from := "/casaOS/server/conf"
		if fds, err = ioutil.ReadDir(from); err != nil {
			return err
		}

		for _, fd := range fds {
			fExt := path.Ext(fd.Name())
			if fExt != ".json" {
				continue
			}
			srcfp := path.Join(from, fd.Name())
			dstfp := path.Join(to, fd.Name())
			source, err := os.Open(srcfp)
			if err != nil {
				return err
			}
			defer source.Close()

			destination, err := os.Create(dstfp)
			if err != nil {
				return err
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				return err
			}
		}

	}

	_logger.Info("update done")
	return nil
}

func (u *migrationTool) PostMigrate() error {
	return nil
}

func NewMigrationToolFor_035() interfaces.MigrationTool {
	return &migrationTool{}
}
