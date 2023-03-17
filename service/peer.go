/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-07-26 18:13:22
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-04 20:10:31
 * @FilePath: /CasaOS/service/connections.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"github.com/IceWhaleTech/CasaOS/service/model"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type PeerService interface {
	GetPeerByUserAgent(ua string) (m model2.PeerDriveDBModel)
	GetPeerByID(id string) (m model2.PeerDriveDBModel)
	GetPeers() (peers []model2.PeerDriveDBModel)
	CreatePeer(m *model2.PeerDriveDBModel)
	DeletePeer(id string)
	GetPeerByName(name string) (m model2.PeerDriveDBModel)
}

type peerStruct struct {
	db *gorm.DB
}

func (s *peerStruct) GetPeerByName(name string) (m model2.PeerDriveDBModel) {
	s.db.Where("display_name = ?", name).First(&m)
	return
}
func (s *peerStruct) GetPeerByUserAgent(ua string) (m model2.PeerDriveDBModel) {
	s.db.Where("user_agent = ?", ua).First(&m)
	return
}
func (s *peerStruct) GetPeerByID(id string) (m model2.PeerDriveDBModel) {
	s.db.Where("id = ?", id).First(&m)
	return
}
func (s *peerStruct) GetPeers() (peers []model2.PeerDriveDBModel) {
	s.db.Order("updated desc").Find(&peers)
	return
}
func (s *peerStruct) CreatePeer(m *model2.PeerDriveDBModel) {

	s.db.Create(m)
}

func (s *peerStruct) DeletePeer(id string) {
	s.db.Where("id= ?", id).Delete(&model.PeerDriveDBModel{})
}

func NewPeerService(db *gorm.DB) PeerService {
	return &peerStruct{db: db}
}
