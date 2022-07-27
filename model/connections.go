/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-27 10:30:43
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-27 11:36:08
 * @FilePath: /CasaOS/model/connections.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type Connections struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Directory  string `json:"directory"`
	MountPoint string `json:"mount_point"`
}
