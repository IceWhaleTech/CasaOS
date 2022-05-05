package service

import (
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type DownRecordService interface {
	AddDownRecord(m model2.PersonDownRecordDBModel)
	GetDownloadListByFrom(id string) []model2.PersonDownRecordDBModel
	GetDownloadListByPath(path string) (list []model2.PersonDownRecordDBModel)
}
type downRecordService struct {
	db *gorm.DB
}

func (d *downRecordService) AddDownRecord(m model2.PersonDownRecordDBModel) {
	d.db.Create(&m)
}

func (d *downRecordService) GetDownloadListByFrom(id string) []model2.PersonDownRecordDBModel {
	var m []model2.PersonDownRecordDBModel
	d.db.Model(m).Where("from = ?", id).Find(&m)
	return m
}
func (d *downRecordService) GetDownloadListByPath(path string) (list []model2.PersonDownRecordDBModel) {
	d.db.Where("path = ?", path).Find(&list)
	return
}

func NewDownRecordService(db *gorm.DB) DownRecordService {
	return &downRecordService{db: db}
}
