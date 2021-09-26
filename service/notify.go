package service

import (
	"gorm.io/gorm"
	"oasis/service/model"
	"oasis/types"
)

type NotifyServer interface {
	GetLog(id string) model.AppNotify
	AddLog(log model.AppNotify)
	UpdateLog(log model.AppNotify)
	DelLog(id string)
	GetList() (list []model.AppNotify)
	MarkRead(id string, state int)
}

type notifyServer struct {
	db *gorm.DB
}

func (i notifyServer) GetList() (list []model.AppNotify) {
	i.db.Where("state=? or state=?", types.NOTIFY_DYNAMICE, types.NOTIFY_UNREAD).Find(&list)
	return
}

func (i *notifyServer) AddLog(log model.AppNotify) {
	i.db.Create(&log)
}

func (i *notifyServer) UpdateLog(log model.AppNotify) {
	i.db.Save(&log)
}

func (i *notifyServer) GetLog(id string) model.AppNotify {
	var log model.AppNotify
	i.db.Where("custom_id = ? ", id).First(&log)
	return log
}
func (i *notifyServer) MarkRead(id string, state int) {
	i.db.Update("state=", state).Where("custom_id = ? ", id)
}
func (i *notifyServer) DelLog(id string) {
	var log model.AppNotify
	i.db.Where("custom_id = ?", id).Delete(&log)
}

func NewNotifyService(db *gorm.DB) NotifyServer {
	return &notifyServer{db: db}
}
