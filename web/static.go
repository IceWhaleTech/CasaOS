/*
 * @Author: Jerryk jerry@icewhale.org
 * @Date: 2022-02-18 10:20:10
 * @LastEditors: Jerryk jerry@icewhale.org
 * @LastEditTime: 2022-05-16 14:56:14
 * @FilePath: \CasaOS-UI\public\static.go
 * @Description:
 *
 * Copyright (c) 2022 by IceWhale, All Rights Reserved.
 */

package web

import "embed"

//go:embed index.html favicon.svg browserconfig.xml site.webmanifest robots.txt img js fonts *.worker.js
var Static embed.FS
