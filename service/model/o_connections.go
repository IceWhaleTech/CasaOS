/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 17:17:57
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-01 17:08:08
 * @FilePath: /CasaOS/service/model/o_connections.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type ConnectionsDBModel struct {
	ID          uint   `gorm:"column:id;primary_key" json:"id"`
	Updated     int64  `gorm:"autoUpdateTime"`
	Created     int64  `gorm:"autoCreateTime"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Status      string `json:"status"`
	Directories string `json:"directories"` // string array
	MountPoint  string `json:"mount_point"` //parent directory of mount point
}

func (p *ConnectionsDBModel) TableName() string {
	return "o_connections"
}
