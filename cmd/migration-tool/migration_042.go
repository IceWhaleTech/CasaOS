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
	"os"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
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

	if minorVersion > 4 {
		return false, nil
	}

	if minorVersion == 4 && patchVersion != 2 {
		return false, nil
	}

	_logger.Info("Migration is needed for a CasaOS version 0.4.2 ")
	return true, nil
}

func (u *migrationTool) PreMigrate() error {
	return nil
}

func (u *migrationTool) Migrate() error {
	_logger.Info("Migration is started for a CasaOS version 0.4.2 ")
	command.OnlyExec("systemctl stop rclone.service")
	os.Remove("/usr/lib/systemd/system/rclone.service")
	command.OnlyExec("systemctl daemon-reload")
	return nil
}

func (u *migrationTool) PostMigrate() error {
	return nil
}

func NewMigrationDummy() interfaces.MigrationTool {
	return &migrationTool{}
}
