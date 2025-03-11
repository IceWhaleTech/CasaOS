/*
* @Author: Drew Fitzgerald/Sheep26 drew@sheepland.xyz
* @Date: 2024-12-03
* @LastEditors: Drew Fitzgerald/Sheep26
* @LastEditTime: 2024-12-03
* @FilePath: /CasaOS/service/model/o_smb_users.go
* @Website: https://www.casaos.io
* Copyright (c) 2024 by icewhale, All Rights Reserved.
*/

package model

type SMBUsers struct {
	Name        string   "json:name"
	Password    string    "json:password"
}