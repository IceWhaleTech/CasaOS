/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-31 17:04:02
 * @FilePath: /CasaOS/pkg/config/config.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package config

import (
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
)

var CasaOSConfigFilePath = filepath.Join(constants.DefaultConfigPath, "casaos.conf")
