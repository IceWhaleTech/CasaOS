/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-27 18:42:42
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-27 18:43:08
 * @FilePath: /CasaOS/model/notify/person.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package notify

type Person struct {
	ShareId string `json:"share_id"`
	Type    string `json:"type"`
}
