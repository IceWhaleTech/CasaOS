package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/IceWhaleTech/CasaOS/pkg/generic_sync"

	"github.com/IceWhaleTech/CasaOS/model"

	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/internal/op"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
)

type StoragesService interface {
	HasStorage(mountPath string) bool
	CreateStorage(ctx context.Context, storage model.Storage) (uint, error)
	LoadStorage(ctx context.Context, storage model.Storage) error
	EnableStorage(ctx context.Context, id uint) error
	DisableStorage(ctx context.Context, id uint) error
	UpdateStorage(ctx context.Context, storage model.Storage) error
	DeleteStorageById(ctx context.Context, id uint) error
	MustSaveDriverStorage(driver driver.Driver) error
	GetStorageVirtualFilesByPath(prefix string) []model.Obj
	initStorage(ctx context.Context, storage model.Storage, storageDriver driver.Driver, setMountPath func(d driver.Driver, ctx context.Context) string) (err error)
	InitStorages()
	GetBalancedStorage(path string) driver.Driver
}

type storagesStruct struct {
}

// Although the driver type is stored,
// there is a storage in each driver,
// so it should actually be a storage, just wrapped by the driver
var storagesMap generic_sync.MapOf[string, driver.Driver]

func GetAllStorages() []driver.Driver {
	return storagesMap.Values()
}

func (s *storagesStruct) HasStorage(mountPath string) bool {
	return storagesMap.Has(utils.FixAndCleanPath(mountPath))
}

func GetStorageByMountPath(mountPath string) (driver.Driver, error) {
	mountPath = utils.FixAndCleanPath(mountPath)
	storageDriver, ok := storagesMap.Load(mountPath)
	if !ok {
		return nil, errors.Errorf("no mount path for an storage is: %s", mountPath)
	}
	return storageDriver, nil
}

// CreateStorage Save the storage to database so storage can get an id
// then instantiate corresponding driver and save it in memory
func (s *storagesStruct) CreateStorage(ctx context.Context, storage model.Storage) (uint, error) {
	storage.Modified = time.Now()
	storage.MountPath = utils.FixAndCleanPath(storage.MountPath)
	var err error
	// check driver first
	driverName := storage.Driver
	driverNew, err := op.GetDriverNew(driverName)
	if err != nil {
		return 0, errors.WithMessage(err, "failed get driver new")
	}
	storageDriver := driverNew()
	// // insert storage to database
	// err = MyService.Storage().CreateStorage(&storage)
	// if err != nil {

	// 	return storage.ID, errors.WithMessage(err, "failed create storage in database")
	// }
	// already has an id
	err = s.initStorage(ctx, storage, storageDriver, func(d driver.Driver, ctx context.Context) string {
		u, _ := d.GetUserInfo(ctx)
		if len(u) > 0 {
			a := strings.Split(u, "@")
			u = a[0]
		}
		return u
	})
	if err != nil {
		s.DeleteStorageById(ctx, storage.ID)
		return storage.ID, errors.Wrap(err, "failed init storage")
	}

	go op.CallStorageHooks("add", storageDriver)

	logger.Error("storage created", zap.Any("storage", storageDriver))
	return storage.ID, nil
}

// LoadStorage load exist storage in db to memory
func (s *storagesStruct) LoadStorage(ctx context.Context, storage model.Storage) error {
	storage.MountPath = utils.FixAndCleanPath(storage.MountPath)
	// check driver first
	driverName := storage.Driver
	driverNew, err := op.GetDriverNew(driverName)
	if err != nil {
		return errors.WithMessage(err, "failed get driver new")
	}
	storageDriver := driverNew()

	err = s.initStorage(ctx, storage, storageDriver, nil)
	go op.CallStorageHooks("add", storageDriver)
	logger.Info("storage created", zap.Any("storage", storageDriver))
	return err
}

// initStorage initialize the driver and store to storagesMap
func (s *storagesStruct) initStorage(ctx context.Context, storage model.Storage, storageDriver driver.Driver, setMountPath func(d driver.Driver, ctx context.Context) string) (err error) {
	storageDriver.SetStorage(storage)
	driverStorage := storageDriver.GetStorage()

	// Unmarshal Addition

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	err = json.UnmarshalFromString(driverStorage.Addition, storageDriver.GetAddition())
	if err == nil {
		err = storageDriver.Init(ctx)
	}
	if setMountPath != nil {
		driverStorage.MountPath += "_" + setMountPath(storageDriver, ctx)

	}
	if s.HasStorage(driverStorage.MountPath) {
		return errors.New("mount path already exists")
	}
	storageDriver.SetStorage(*driverStorage)
	storagesMap.Store(driverStorage.MountPath, storageDriver)

	if err != nil {
		driverStorage.SetStatus(err.Error())
		err = errors.Wrap(err, "failed init storage")
	} else {
		driverStorage.SetStatus(op.WORK)
	}

	err = s.MustSaveDriverStorage(storageDriver)

	return err
}

func (s *storagesStruct) EnableStorage(ctx context.Context, id uint) error {
	// storage, err := MyService.Storage().GetStorageById(id)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get storage")
	// }
	// if !storage.Disabled {
	// 	return errors.Errorf("this storage have enabled")
	// }
	// storage.Disabled = false
	// err = MyService.Storage().UpdateStorage(storage)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed update storage in db")
	// }
	// err = s.LoadStorage(ctx, *storage)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed load storage")
	// }
	return nil
}

func (s *storagesStruct) DisableStorage(ctx context.Context, id uint) error {
	// storage, err := MyService.Storage().GetStorageById(id)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get storage")
	// }
	// if storage.Disabled {
	// 	return errors.Errorf("this storage have disabled")
	// }
	// storageDriver, err := GetStorageByMountPath(storage.MountPath)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get storage driver")
	// }
	// // drop the storage in the driver
	// if err := storageDriver.Drop(ctx); err != nil {
	// 	return errors.Wrap(err, "failed drop storage")
	// }
	// // delete the storage in the memory
	// storage.Disabled = true
	// err = MyService.Storage().UpdateStorage(storage)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed update storage in db")
	// }
	// storagesMap.Delete(storage.MountPath)
	// go op.CallStorageHooks("del", storageDriver)
	return nil
}

// UpdateStorage update storage
// get old storage first
// drop the storage then reinitialize
func (s *storagesStruct) UpdateStorage(ctx context.Context, storage model.Storage) error {
	// oldStorage, err := MyService.Storage().GetStorageById(storage.ID)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get old storage")
	// }
	// if oldStorage.Driver != storage.Driver {
	// 	return errors.Errorf("driver cannot be changed")
	// }
	// storage.Modified = time.Now()
	// storage.MountPath = utils.FixAndCleanPath(storage.MountPath)
	// err = MyService.Storage().UpdateStorage(&storage)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed update storage in database")
	// }
	// if storage.Disabled {
	// 	return nil
	// }
	// storageDriver, err := GetStorageByMountPath(oldStorage.MountPath)
	// if oldStorage.MountPath != storage.MountPath {
	// 	// mount path renamed, need to drop the storage
	// 	storagesMap.Delete(oldStorage.MountPath)
	// }
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get storage driver")
	// }
	// err = storageDriver.Drop(ctx)
	// if err != nil {
	// 	return errors.Wrapf(err, "failed drop storage")
	// }

	// err = s.initStorage(ctx, storage, storageDriver, nil)
	// go op.CallStorageHooks("update", storageDriver)

	// logger.Info("storage updated", zap.Any("storage", storageDriver))
	//return err
	return nil
}

func (s *storagesStruct) DeleteStorageById(ctx context.Context, id uint) error {
	// storage, err := MyService.Storage().GetStorageById(id)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed get storage")
	// }
	// if !storage.Disabled {
	// 	storageDriver, err := GetStorageByMountPath(storage.MountPath)
	// 	if err == nil {
	// 		// drop the storage in the driver
	// 		if err := storageDriver.Drop(ctx); err != nil {
	// 			return errors.Wrapf(err, "failed drop storage")
	// 		}
	// 		// delete the storage in the memory
	// 		storagesMap.Delete(storage.MountPath)
	// 	}

	// 	go op.CallStorageHooks("del", storageDriver)
	// }
	// // delete the storage in the database
	// if err := MyService.Storage().DeleteStorageById(id); err != nil {
	// 	return errors.WithMessage(err, "failed delete storage in database")
	// }
	return nil
}

// MustSaveDriverStorage call from specific driver
func (s *storagesStruct) MustSaveDriverStorage(driver driver.Driver) error {
	err := saveDriverStorage(driver)
	if err != nil {
		logger.Error("failed save driver storage", zap.Any("err", err))
	}
	return err
}

func saveDriverStorage(driver driver.Driver) error {
	// storage := driver.GetStorage()
	// addition := driver.GetAddition()

	// var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// str, err := json.MarshalToString(addition)
	// if err != nil {
	// 	return errors.Wrap(err, "error while marshal addition")
	// }
	// storage.Addition = str
	// err = MyService.Storage().UpdateStorage(storage)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed update storage in database")
	// }
	return nil
}

// getStoragesByPath get storage by longest match path, contains balance storage.
// for example, there is /a/b,/a/c,/a/d/e,/a/d/e.balance
// getStoragesByPath(/a/d/e/f) => /a/d/e,/a/d/e.balance
func getStoragesByPath(path string) []driver.Driver {
	storages := make([]driver.Driver, 0)
	curSlashCount := 0
	storagesMap.Range(func(mountPath string, value driver.Driver) bool {
		mountPath = utils.GetActualMountPath(mountPath)
		// is this path
		if utils.IsSubPath(mountPath, path) {
			slashCount := strings.Count(utils.PathAddSeparatorSuffix(mountPath), "/")
			// not the longest match
			if slashCount > curSlashCount {
				storages = storages[:0]
				curSlashCount = slashCount
			}
			if slashCount == curSlashCount {
				storages = append(storages, value)
			}
		}
		return true
	})
	// make sure the order is the same for same input
	sort.Slice(storages, func(i, j int) bool {
		return storages[i].GetStorage().MountPath < storages[j].GetStorage().MountPath
	})
	return storages
}

// GetStorageVirtualFilesByPath Obtain the virtual file generated by the storage according to the path
// for example, there are: /a/b,/a/c,/a/d/e,/a/b.balance1,/av
// GetStorageVirtualFilesByPath(/a) => b,c,d
func (s *storagesStruct) GetStorageVirtualFilesByPath(prefix string) []model.Obj {
	files := make([]model.Obj, 0)
	storages := storagesMap.Values()
	sort.Slice(storages, func(i, j int) bool {
		if storages[i].GetStorage().Order == storages[j].GetStorage().Order {
			return storages[i].GetStorage().MountPath < storages[j].GetStorage().MountPath
		}
		return storages[i].GetStorage().Order < storages[j].GetStorage().Order
	})

	prefix = utils.FixAndCleanPath(prefix)
	set := mapset.NewSet[string]()
	for _, v := range storages {
		mountPath := utils.GetActualMountPath(v.GetStorage().MountPath)
		// Exclude prefix itself and non prefix
		if len(prefix) >= len(mountPath) || !utils.IsSubPath(prefix, mountPath) {
			continue
		}
		name := strings.SplitN(strings.TrimPrefix(mountPath[len(prefix):], "/"), "/", 2)[0]
		if set.Add(name) {
			files = append(files, &model.Object{
				Name:     name,
				Size:     0,
				Modified: v.GetStorage().Modified,
				IsFolder: true,
			})
		}
	}
	return files
}

var balanceMap generic_sync.MapOf[string, int]

// GetBalancedStorage get storage by path
func (s *storagesStruct) GetBalancedStorage(path string) driver.Driver {
	path = utils.FixAndCleanPath(path)
	storages := getStoragesByPath(path)
	storageNum := len(storages)
	switch storageNum {
	case 0:
		return nil
	case 1:
		return storages[0]
	default:
		virtualPath := utils.GetActualMountPath(storages[0].GetStorage().MountPath)
		i, _ := balanceMap.LoadOrStore(virtualPath, 0)
		i = (i + 1) % storageNum
		balanceMap.Store(virtualPath, i)
		return storages[i]
	}
}
func (s *storagesStruct) InitStorages() {
	// storages, err := MyService.Storage().GetEnabledStorages()
	// if err != nil {
	// 	logger.Error("failed get enabled storages", zap.Any("err", err))
	// }
	// go func(storages []model.Storage) {
	// 	for i := range storages {
	// 		err := s.LoadStorage(context.Background(), storages[i])
	// 		if err != nil {
	// 			logger.Error("failed get enabled storages", zap.Any("err", err))
	// 		} else {
	// 			logger.Info("success load storage", zap.String("mount_path", storages[i].MountPath))
	// 		}
	// 	}
	// 	conf.StoragesLoaded = true
	// }(storages)

}
func NewStoragesService() StoragesService {
	return &storagesStruct{}
}
