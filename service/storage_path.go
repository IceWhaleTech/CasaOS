package service

import (
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StoragePathService interface {
	GetStorageAndActualPath(rawPath string) (storage driver.Driver, actualPath string, err error)
}

type storagePathStruct struct {
}

func (s *storagePathStruct) GetStorageAndActualPath(rawPath string) (storage driver.Driver, actualPath string, err error) {
	rawPath = utils.FixAndCleanPath(rawPath)
	storage = MyService.Storages().GetBalancedStorage(rawPath)
	if storage == nil {
		err = errors.Errorf("can't find storage with rawPath: %s", rawPath)
		return
	}
	logger.Info("use storage", zap.Any("storage mount path", storage.GetStorage().MountPath))
	mountPath := utils.GetActualMountPath(storage.GetStorage().MountPath)
	actualPath = utils.FixAndCleanPath(strings.TrimPrefix(rawPath, mountPath))
	return
}
func NewStoragePathService() StoragePathService {
	return &storagePathStruct{}
}
