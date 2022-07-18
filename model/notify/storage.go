/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-15 10:43:00
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-15 10:56:17
 * @FilePath: /CasaOS/model/notify/storage.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package notify

type StorageMessage struct {
	Type   string `json:"type"`   //sata,usb
	Action string `json:"action"` //remove add
	Path   string `json:"path"`
	Volume string `json:"volume"`
	Size   uint64 `json:"size"`
}
