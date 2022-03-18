package service

import (
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type DownloadService interface {
	AddDownloadTask(m model2.PersionDownloadDBModel)   //添加下载任务
	EditDownloadState(m model2.PersionDownloadDBModel) //只修改状态
	SaveDownload(m model2.PersionDownloadDBModel)
	DelDownload(uuid string)
	GetDownloadById(uuid string) model2.PersionDownloadDBModel
	GetDownloadListByState(state string) []model2.PersionDownloadDBModel
	SetDownloadError(m model2.PersionDownloadDBModel)
}
type downloadService struct {
	db *gorm.DB
}

func (d *downloadService) AddDownloadTask(m model2.PersionDownloadDBModel) {
	d.db.Create(&m)
}
func (d *downloadService) EditDownloadState(m model2.PersionDownloadDBModel) {
	d.db.Model(&m).Where("uuid = ?", m.UUID).Update("state", m.State)
}

//failed during download
func (d *downloadService) SetDownloadError(m model2.PersionDownloadDBModel) {
	d.db.Model(&m).Updates(m)
}

func (d *downloadService) DelDownload(uuid string) {
	var m model2.PersionDownloadDBModel
	d.db.Where("uuid = ?", uuid).Delete(&m)
}
func (d *downloadService) GetDownloadById(uuid string) model2.PersionDownloadDBModel {
	var m model2.PersionDownloadDBModel
	d.db.Model(m).Where("uuid = ?", uuid).First(&m)
	return m
}
func (d *downloadService) GetDownloadListByState(state string) (list []model2.PersionDownloadDBModel) {
	if len(state) == 0 {
		d.db.Find(&list)

	} else {
		d.db.Where("state = ?", state).Find(&list)
	}

	return
}

func (d *downloadService) SaveDownload(m model2.PersionDownloadDBModel) {
	d.db.Save(&m)
}
func NewDownloadService(db *gorm.DB) DownloadService {
	return &downloadService{db: db}
}
