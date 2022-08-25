/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-08-23 18:09:11
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-24 18:02:59
 * @FilePath: /CasaOS/cmd/migration-tool/main.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package main

import (
	"flag"
	"fmt"
	"os"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

const (
	casaosServiceName = "casaos.service"
)

func main() {
	versionFlag := flag.Bool("v", false, "version")
	debugFlag := flag.Bool("d", true, "debug")
	forceFlag := flag.Bool("f", false, "force")
	flag.Parse()

	if *versionFlag {
		fmt.Println(common.Version)
		os.Exit(0)
	}

	if os.Getuid() != 0 {
		os.Exit(1)
	}

	if *debugFlag {
		//		_logger.DebugMode = true
	}

	if !*forceFlag {
		serviceEnabled, err := systemctl.IsServiceEnabled(casaosServiceName)
		if err != nil {
			panic(err)
		}

		if serviceEnabled {
			//_logger.Info("%s is already enabled. If migration is still needed, try with -f.", userServiceName)
			os.Exit(1)
		}
	}

	migrationTools := []interfaces.MigrationTool{
		NewMigrationToolFor_035(),
	}

	var selectedMigrationTool interfaces.MigrationTool

	// look for the right migration tool matching current version
	for _, tool := range migrationTools {
		migrationNeeded, err := tool.IsMigrationNeeded()
		if err != nil {
			panic(err)
		}

		if migrationNeeded {
			selectedMigrationTool = tool
			break
		}
	}

	if selectedMigrationTool == nil {
		return
	}

	if err := selectedMigrationTool.PreMigrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.Migrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.PostMigrate(); err != nil {
		panic(err)
	}
}
