/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-23 17:28:51
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-23 17:28:54
 * @FilePath: /CasaOS/web/static.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package web

import "embed"

//go:embed index.html favicon.svg browserconfig.xml site.webmanifest robots.txt img js fonts css
var Static embed.FS
