/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 17:17:57
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-27 10:44:46
 * @FilePath: /CasaOS/service/model/o_connections.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type ConnectionsDBModel struct {
	ID         uint   `gorm:"column:id;primary_key" json:"id"`
	Updated    int64  `gorm:"autoUpdateTime"`
	Created    int64  `gorm:"autoCreateTime"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Directory  string `json:"directory"`
	MountPoint string `json:"mount_point"`
	Status     string `json:"status"`
}

func (p *ConnectionsDBModel) TableName() string {
	return "o_connections"
}
