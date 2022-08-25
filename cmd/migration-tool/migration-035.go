/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-08-24 17:36:00
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-24 18:11:16
 * @FilePath: /CasaOS/cmd/migration-tool/migration-035.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package main

import (
	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
)

type migrationTool struct{}

func (u *migrationTool) IsMigrationNeeded() (bool, error) {

	minorVersion, err := version.DetectMinorVersion()
	if err != nil {
		return false, err
	}

	if minorVersion != 3 {
		return false, nil
	}

	// this is the best way to tell if CasaOS version is between 0.3.3 and 0.3.5

	//return true, nil
	return false, nil
}

func (u *migrationTool) PreMigrate() error {

	return nil
}

func (u *migrationTool) Migrate() error {
	return nil
}

func (u *migrationTool) PostMigrate() error {
	return nil
}

func NewMigrationToolFor_035() interfaces.MigrationTool {
	return &migrationTool{}
}
