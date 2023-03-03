package service

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/pkg/mount"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/cmd/mountlib"
	"github.com/rclone/rclone/fs"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/rc"
	"github.com/rclone/rclone/vfs/vfscommon"
	"go.uber.org/zap"
)

type StorageService interface {
	MountStorage(mountPoint, fs string) error
	UnmountStorage(mountPoint string) error
	UnmountAllStorage()
	GetStorages() (httper.MountList, error)
	CreateConfig(data rc.Params, name string, t string) error
	CheckAndMountByName(name string) error
	CheckAndMountAll() error
	GetConfigByName(name string) []string
	GetAttributeValueByName(name, key string) string
	DeleteConfigByName(name string)
	GetConfig() (httper.RemotesResult, error)
}

type storageStruct struct {
}

var MountLists map[string]*mountlib.MountPoint
var mountMu sync.Mutex

func (s *storageStruct) MountStorage(mountPoint, deviceName string) error {
	file.IsNotExistMkDir(mountPoint)
	mountMu.Lock()
	defer mountMu.Unlock()
	currentFS, err := fs.NewFs(context.TODO(), deviceName+":")
	if err != nil {
		logger.Error("when CheckAndMountAll then", zap.Error(err))
		return err
	}
	mountOptin := mountlib.Options{
		MaxReadAhead:  128 * 1024,
		AttrTimeout:   1 * time.Second,
		DaemonWait:    60 * time.Second,
		NoAppleDouble: true,
		NoAppleXattr:  false,
		AsyncRead:     true,
		AllowOther:    true,
	}
	vfsOpt := vfscommon.Options{
		NoModTime:          false,
		NoChecksum:         false,
		NoSeek:             false,
		DirCacheTime:       5 * 60 * time.Second,
		PollInterval:       time.Minute,
		ReadOnly:           false,
		Umask:              0,
		UID:                ^uint32(0), // these values instruct WinFSP-FUSE to use the current user
		GID:                ^uint32(0), // overridden for non windows in mount_unix.go
		DirPerms:           os.FileMode(0777),
		FilePerms:          os.FileMode(0666),
		CacheMode:          3,
		CacheMaxAge:        3600 * time.Second,
		CachePollInterval:  60 * time.Second,
		ChunkSize:          128 * fs.Mebi,
		ChunkSizeLimit:     -1,
		CacheMaxSize:       -1,
		CaseInsensitive:    runtime.GOOS == "windows" || runtime.GOOS == "darwin", // default to true on Windows and Mac, false otherwise
		WriteWait:          1000 * time.Millisecond,
		ReadWait:           20 * time.Millisecond,
		WriteBack:          5 * time.Second,
		ReadAhead:          0 * fs.Mebi,
		UsedIsSize:         false,
		DiskSpaceTotalSize: -1,
	}

	mnt := mountlib.NewMountPoint(mount.MountFn, mountPoint, currentFS, &mountOptin, &vfsOpt)
	_, err = mnt.Mount()
	if err != nil {
		logger.Error("when CheckAndMountAll then", zap.Error(err))
		return err
	}
	go func() {
		if err = mnt.Wait(); err != nil {
			log.Printf("unmount FAILED: %v", err)
			return
		}
		mountMu.Lock()
		defer mountMu.Unlock()
		delete(MountLists, mountPoint)
	}()
	return nil
}
func (s *storageStruct) UnmountStorage(mountPoint string) error {

	err := MountLists[mountPoint].Unmount()
	if err != nil {
		logger.Error("when umount then", zap.Error(err))
		return err
	}
	return nil
}
func (s *storageStruct) UnmountAllStorage() {
	for _, v := range MountLists {
		err := v.Unmount()
		if err != nil {
			logger.Error("when umount then", zap.Error(err))
		}
	}
}
func (s *storageStruct) GetStorages() (httper.MountList, error) {
	ls := httper.MountList{}
	list := []httper.MountPoints{}
	for _, v := range MountLists {
		list = append(list, httper.MountPoints{
			MountPoint: v.MountPoint,
			Fs:         v.Fs.Name(),
		})
	}
	ls.MountPoints = list
	return ls, nil
	// return httper.GetMountList()
}
func (s *storageStruct) CreateConfig(data rc.Params, name string, t string) error {
	_, err := rconfig.CreateRemote(context.Background(), name, t, data, rconfig.UpdateRemoteOpt{State: "*oauth-islocal,teamdrive,,", NonInteractive: true})
	return err
}
func (s *storageStruct) CheckAndMountByName(name string) error {

	mountPoint, found := rconfig.LoadedData().GetValue(name, "mount_point")
	if !found && len(mountPoint) == 0 {
		logger.Error("when CheckAndMountAll then mountpint is empty", zap.String("mountPoint", mountPoint), zap.String("fs", name))
		return errors.New("mountpoint is empty")
	}
	return MyService.Storage().MountStorage(mountPoint, name)
}

func (s *storageStruct) CheckAndMountAll() error {
	section := rconfig.LoadedData().GetSectionList()

	logger.Info("when CheckAndMountAll section", zap.Any("section", section))
	for _, v := range section {
		command.OnlyExec("umount /mnt/" + v)
		mountPoint, found := rconfig.LoadedData().GetValue(v, "mount_point")

		if !found && len(mountPoint) == 0 {
			logger.Info("when CheckAndMountAll then mountpint is empty", zap.String("mountPoint", mountPoint), zap.String("fs", v))
			continue
		}
		err := MyService.Storage().MountStorage(mountPoint, v)
		if err != nil {
			logger.Error("when CheckAndMountAll then", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *storageStruct) GetConfigByName(name string) []string {
	return rconfig.LoadedData().GetKeyList(name)
}

func (s *storageStruct) GetAttributeValueByName(name, key string) string {
	value, found := rconfig.LoadedData().GetValue(name, key)
	if !found {
		return ""
	}
	return value
}

func (s *storageStruct) DeleteConfigByName(name string) {
	rconfig.DeleteRemote(name)
}
func (s *storageStruct) GetConfig() (httper.RemotesResult, error) {
	//TODO: check data
	// section, err := httper.GetAllConfigName()
	// if err != nil {
	// 	return httper.RemotesResult{}, err
	// }
	// return section, nil
	return httper.RemotesResult{}, nil
}
func NewStorageService() StorageService {
	return &storageStruct{}
}
