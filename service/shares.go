/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 11:21:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-27 15:27:38
 * @FilePath: /CasaOS/service/shares.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service/model"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type SharesService interface {
	GetSharesList() (shares []model2.SharesDBModel)
	GetSharesByPath(path string) (shares []model2.SharesDBModel)
	GetSharesByName(name string) (shares []model2.SharesDBModel)
	CreateShare(share model2.SharesDBModel)
	DeleteShare(id string)
	UpdateConfigFile()
}

type sharesStruct struct {
	db *gorm.DB
}

func (s *sharesStruct) GetSharesByName(name string) (shares []model2.SharesDBModel) {
	s.db.Select("title,anonymous,path,id").Where("name = ?", name).Find(&shares)
	return
}
func (s *sharesStruct) GetSharesByPath(path string) (shares []model2.SharesDBModel) {
	s.db.Select("title,anonymous,path,id").Where("path = ?", path).Find(&shares)
	return
}
func (s *sharesStruct) GetSharesList() (shares []model2.SharesDBModel) {
	s.db.Select("title,anonymous,path,id").Find(&shares)
	return
}
func (s *sharesStruct) CreateShare(share model2.SharesDBModel) {
	s.db.Create(&share)
	s.UpdateConfigFile()
}
func (s *sharesStruct) DeleteShare(id string) {
	s.db.Where("id= ?", id).Delete(&model.SharesDBModel{})
	s.UpdateConfigFile()
}

func (s *sharesStruct) UpdateConfigFile() {
	shares := []model2.SharesDBModel{}
	s.db.Select("title,anonymous,path").Find(&shares)
	//generated config file
	var configStr = ""
	for _, share := range shares {
		dirName := filepath.Base(share.Path)
		configStr += `
[` + dirName + `]
public = Yes
path = ` + share.Path + `
browseable = Yes
read only = No
guest ok = Yes
create mask = 0777
directory mask = 0777

`
	}
	//write config file
	file.WriteToPath([]byte(configStr), "/etc/samba", "smb.casa.conf")
	//restart samba
	command2.ExecResultStrArray("source " + config.AppInfo.ShellPath + "/helper.sh ;RestartSMBD")
}
func NewSharesService(db *gorm.DB) SharesService {
	return &sharesStruct{db: db}
}
