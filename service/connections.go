/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 18:13:22
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-27 13:44:12
 * @FilePath: /CasaOS/service/connections.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/service/model"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type ConnectionsService interface {
	GetConnectionsList() (connections []model2.ConnectionsDBModel)
	GetConnectionByDirectory(directory string) (connections []model2.ConnectionsDBModel)
	GetConnectionByID(id string) (connections model2.ConnectionsDBModel)
	CreateConnection(connection *model2.ConnectionsDBModel)
	DeleteConnection(id string)
	MountSmaba(connection *model2.ConnectionsDBModel) string
	UnmountSmaba(mountPoint string) string
}

type connectionsStruct struct {
	db *gorm.DB
}

func (s *connectionsStruct) GetConnectionByDirectory(directory string) (connections []model2.ConnectionsDBModel) {
	s.db.Select("username,host,directory,status,mount_point,id").Where("directory = ?", directory).Find(&connections)
	return
}
func (s *connectionsStruct) GetConnectionByID(id string) (connections model2.ConnectionsDBModel) {
	s.db.Select("username,password,host,directory,status,mount_point,id").Where("id = ?", id).First(&connections)
	return
}
func (s *connectionsStruct) GetConnectionsList() (connections []model2.ConnectionsDBModel) {
	s.db.Select("username,host,port,directory,status,mount_point,id").Find(&connections)
	return
}
func (s *connectionsStruct) CreateConnection(connection *model2.ConnectionsDBModel) {
	s.db.Create(connection)
}
func (s *connectionsStruct) DeleteConnection(id string) {
	s.db.Where("id= ?", id).Delete(&model.ConnectionsDBModel{})
}

func (s *connectionsStruct) MountSmaba(connection *model2.ConnectionsDBModel) string {
	str := command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;MountCIFS " + connection.Username + " " + connection.Host + " " + connection.Directory + " " + connection.Port + " " + connection.MountPoint + " " + connection.Password)
	return str
}
func (s *connectionsStruct) UnmountSmaba(mountPoint string) string {
	str := command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;UMountPorintAndRemoveDir " + mountPoint)
	return str
}

func NewConnectionsService(db *gorm.DB) ConnectionsService {
	return &connectionsStruct{db: db}
}
