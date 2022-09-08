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
	"strings"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
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

	_logger.Info("update done")
	return nil
}

func (u *migrationTool) PostMigrate() error {
	return nil
}

func NewMigrationToolFor_036() interfaces.MigrationTool {
	return &migrationTool{}
}
