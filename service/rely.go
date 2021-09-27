package service

import (
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type RelyService interface {
	Create(rely model2.RelyDBModel)
	Delete(id string)
	GetInfo(id string) model2.RelyDBModel
}

type relyService struct {
	db  *gorm.DB
	log loger2.OLog
}

func (r *relyService) Create(rely model2.RelyDBModel) {

	r.db.Create(&rely)

}

//获取我的应用列表
func (r *relyService) GetInfo(id string) model2.RelyDBModel {
	var m model2.RelyDBModel
	r.db.Where("custom_id = ?", id).First(&m)
	return m
}

func (r *relyService) Delete(id string) {
	var c model2.RelyDBModel
	r.db.Where("custom_id = ?", id).Delete(&c)
}

func NewRelyService(db *gorm.DB, log loger2.OLog) RelyService {
	return &relyService{db: db, log: log}
}
