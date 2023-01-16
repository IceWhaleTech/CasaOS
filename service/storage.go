package service

import (
	"fmt"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StorageService interface {
	CreateStorage(storage *model.Storage) error
	UpdateStorage(storage *model.Storage) error
	DeleteStorageById(id uint) error
	GetStorages(pageIndex, pageSize int) ([]model.Storage, int64, error)
	GetStorageById(id uint) (*model.Storage, error)
	GetEnabledStorages() ([]model.Storage, error)
}

type storageStruct struct {
	db *gorm.DB
}

// CreateStorage just insert storage to database
func (s *storageStruct) CreateStorage(storage *model.Storage) error {
	return errors.WithStack(s.db.Create(storage).Error)
}

// UpdateStorage just update storage in database
func (s *storageStruct) UpdateStorage(storage *model.Storage) error {
	return errors.WithStack(s.db.Save(storage).Error)
}

// DeleteStorageById just delete storage from database by id
func (s *storageStruct) DeleteStorageById(id uint) error {
	return errors.WithStack(s.db.Delete(&model.Storage{}, id).Error)
}

// GetStorages Get all storages from database order by index
func (s *storageStruct) GetStorages(pageIndex, pageSize int) ([]model.Storage, int64, error) {
	storageDB := s.db.Model(&model.Storage{})
	var count int64
	if err := storageDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get storages count")
	}
	var storages []model.Storage
	if err := storageDB.Order("`order`").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&storages).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return storages, count, nil
}

// GetStorageById Get Storage by id, used to update storage usually
func (s *storageStruct) GetStorageById(id uint) (*model.Storage, error) {
	var storage model.Storage
	storage.ID = id
	if err := s.db.First(&storage).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &storage, nil
}

func (s *storageStruct) GetEnabledStorages() ([]model.Storage, error) {
	var storages []model.Storage
	if err := s.db.Where(fmt.Sprintf("%s = ?", "disabled"), false).Find(&storages).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return storages, nil
}

func NewStorageService(db *gorm.DB) StorageService {
	return &storageStruct{db: db}
}
