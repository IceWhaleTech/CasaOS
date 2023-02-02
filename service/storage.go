package service

import (
	"io/ioutil"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
)

type StorageService interface {
	MountStorage(mountPoint, fs string) error
	UnmountStorage(mountPoint string) error
	GetStorages() (httper.MountList, error)
	CreateConfig(data map[string]string, name string, t string) error
	CheckAndMountByName(name string) error
	CheckAndMountAll() error
	GetConfigByName(name string) (map[string]string, error)
	DeleteConfigByName(name string) error
}

type storageStruct struct {
}

func (s *storageStruct) MountStorage(mountPoint, fs string) error {
	file.IsNotExistMkDir(mountPoint)
	httper.Mount(mountPoint, fs)
	return nil
}
func (s *storageStruct) UnmountStorage(mountPoint string) error {
	err := httper.Unmount(mountPoint)
	if err == nil {
		dir, _ := ioutil.ReadDir(mountPoint)

		if len(dir) == 0 {
			file.RMDir(mountPoint)
		}
		return nil
	}
	return err
}
func (s *storageStruct) GetStorages() (httper.MountList, error) {
	return httper.GetMountList()
}
func (s *storageStruct) CreateConfig(data map[string]string, name string, t string) error {
	httper.CreateConfig(data, name, t)
	return nil
}
func (s *storageStruct) CheckAndMountByName(name string) error {
	storages, _ := MyService.Storage().GetStorages()
	currentRemote, _ := httper.GetConfigByName(name)
	mountPoint := currentRemote["mount_point"]
	isMount := false
	for _, v := range storages.MountPoints {
		if v.MountPoint == mountPoint {
			isMount = true
			break
		}
	}
	if !isMount {
		MyService.Storage().MountStorage(mountPoint, name+":")
	}
	return nil
}
func (s *storageStruct) CheckAndMountAll() error {
	storages, err := MyService.Storage().GetStorages()
	if err != nil {
		return err
	}
	section, err := httper.GetAllConfigName()
	if err != nil {
		return err
	}
	for _, v := range section.Remotes {
		currentRemote, _ := httper.GetConfigByName(v)
		mountPoint := currentRemote["mount_point"]
		if len(mountPoint) == 0 {
			continue
		}
		isMount := false
		for _, v := range storages.MountPoints {
			if v.MountPoint == mountPoint {
				isMount = true
				break
			}
		}
		if !isMount {
			return MyService.Storage().MountStorage(mountPoint, v+":")
		}
	}
	return nil
}
func (s *storageStruct) GetConfigByName(name string) (map[string]string, error) {
	return httper.GetConfigByName(name)
}
func (s *storageStruct) DeleteConfigByName(name string) error {
	return httper.DeleteConfigByName(name)
}
func NewStorageService() StorageService {
	return &storageStruct{}
}
