package service

import (
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type ShortcutsService interface {
	DeleteData(id string)
	AddData(m model2.ShortcutsDBModel)
	EditData(m model2.ShortcutsDBModel)
	GetList() (list []model2.ShortcutsDBModel)
}
type shortcutsService struct {
	db *gorm.DB
}

func (s *shortcutsService) AddData(m model2.ShortcutsDBModel) {
	s.db.Create(&m)
}
func (s *shortcutsService) EditData(m model2.ShortcutsDBModel) {
	s.db.Save(&m)
}
func (s *shortcutsService) DeleteData(id string) {
	var m model2.ShortcutsDBModel
	s.db.Where("id=?", id).Delete(&m)
}
func (s *shortcutsService) GetList() (list []model2.ShortcutsDBModel) {
	s.db.Order("sort desc,id").Find(&list)
	return list
}
func NewShortcutsService(db *gorm.DB) ShortcutsService {
	return &shortcutsService{db: db}
}
