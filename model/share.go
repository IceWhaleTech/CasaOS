/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 11:12:12
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-27 14:58:55
 * @FilePath: /CasaOS/model/share.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type Shares struct {
	ID          uint     `json:"id"`
	Anonymous   bool     `json:"anonymous"`
	Path        string   `json:"path"`
	Valid_users []string `json:"valid_users"`
}
