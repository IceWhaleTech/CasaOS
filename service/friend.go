package service

import (
	"reflect"

	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type FriendService interface {
	AddFriend(m model2.FriendModel)
	DeleteFriend(m model2.FriendModel)
	EditFriendMark(m model2.FriendModel)
	EditFriendBlock(m model2.FriendModel)
	GetFriendById(m model2.FriendModel) model2.FriendModel
	GetFriendList() (list []model2.FriendModel)
	UpdateAddFriendType(m model2.FriendModel)
	UpdateOrCreate(m model2.FriendModel)
}

type friendService struct {
	db *gorm.DB
}

func (p *friendService) AddFriend(m model2.FriendModel) {
	p.db.Create(&m)
}
func (p *friendService) DeleteFriend(m model2.FriendModel) {
	p.db.Where("token = ?", m.Token).Delete(&m)
}
func (p *friendService) EditFriendMark(m model2.FriendModel) {
	p.db.Model(&m).Where("token = ?", m.Token).Update("mark", m.Mark)
}
func (p *friendService) EditFriendBlock(m model2.FriendModel) {
	p.db.Model(&m).Where("token = ?", m.Token).Update("block", m.Block)
}
func (p *friendService) GetFriendById(m model2.FriendModel) model2.FriendModel {
	p.db.Model(m).Where("token = ?", m.Token).First(&m)
	return m
}

func (p *friendService) GetFriendList() (list []model2.FriendModel) {
	p.db.Select("nick_name", "avatar", "profile", "token", "state", "mark", "block", "version").Find(&list)
	return list
}

func (p *friendService) UpdateOrCreate(m model2.FriendModel) {
	friend := model2.FriendModel{}
	p.db.Where("token = ?", m.Token).First(&friend)
	if reflect.DeepEqual(friend, model2.FriendModel{}) {
		p.db.Create(&m)
	} else {
		p.db.Model(&m).Updates(m)
	}

}

func (p *friendService) UpdateAddFriendType(m model2.FriendModel) {
	p.db.Model(&m).Updates(m)
}

func NewFriendService(db *gorm.DB) FriendService {
	return &friendService{db: db}
}
