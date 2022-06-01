/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-26 14:39:22
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-26 19:08:52
 * @FilePath: /CasaOS/model/notify/message.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package notify

import (
	f "github.com/ambelovsky/gosf"
)

type Message struct {
	Path string    `json:"path"`
	Msg  f.Message `json:"msg"`
}
