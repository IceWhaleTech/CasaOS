package service

import (
	"context"
	"fmt"
	"io/ioutil"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/cmd/mountlib"
	"github.com/rclone/rclone/fs"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/log"
	"github.com/rclone/rclone/vfs"
	"github.com/rclone/rclone/vfs/vfscommon"
	"go.uber.org/zap"
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
	GetConfig() (httper.RemotesResult, error)
}

type storageStruct struct {
}

func (s *storageStruct) MountStorage(mountPoint, fs string) error {
	file.IsNotExistMkDir(mountPoint)
	return httper.Mount(mountPoint, fs)
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
		return MyService.Storage().MountStorage(mountPoint, name+":")
	}
	return nil
}

// mountOptions configures the options from the command line flags
func mountOptions(VFS *vfs.VFS, device string, opt *mountlib.Options) (options []fuse.MountOption) {
	options = []fuse.MountOption{
		fuse.MaxReadahead(uint32(opt.MaxReadAhead)),
		fuse.Subtype("rclone"),
		fuse.FSName(device),

		// Options from benchmarking in the fuse module
		//fuse.MaxReadahead(64 * 1024 * 1024),
		//fuse.WritebackCache(),
	}
	if opt.AsyncRead {
		options = append(options, fuse.AsyncRead())
	}
	if opt.AllowNonEmpty {
		options = append(options, fuse.AllowNonEmptyMount())
	}
	if opt.AllowOther {
		options = append(options, fuse.AllowOther())
	}
	if opt.AllowRoot {
		// options = append(options, fuse.AllowRoot())
		fs.Errorf(nil, "Ignoring --allow-root. Support has been removed upstream - see https://github.com/bazil/fuse/issues/144 for more info")
	}
	if opt.DefaultPermissions {
		options = append(options, fuse.DefaultPermissions())
	}
	if VFS.Opt.ReadOnly {
		options = append(options, fuse.ReadOnly())
	}
	if opt.WritebackCache {
		options = append(options, fuse.WritebackCache())
	}
	if opt.DaemonTimeout != 0 {
		options = append(options, fuse.DaemonTimeout(fmt.Sprint(int(opt.DaemonTimeout.Seconds()))))
	}
	if len(opt.ExtraOptions) > 0 {
		fs.Errorf(nil, "-o/--option not supported with this FUSE backend")
	}
	if len(opt.ExtraFlags) > 0 {
		fs.Errorf(nil, "--fuse-flag not supported with this FUSE backend")
	}
	return options
}

type FS struct {
	*vfs.VFS
	f      fs.Fs
	opt    *mountlib.Options
	server *fusefs.Server
}

func NewFS(VFS *vfs.VFS, opt *mountlib.Options) *FS {
	fsys := &FS{
		VFS: VFS,
		f:   VFS.Fs(),
		opt: opt,
	}
	return fsys
}

// Root returns the root node
func (f *FS) Root() (node fusefs.Node, err error) {
	defer log.Trace("", "")("node=%+v, err=%v", &node, &err)
	root, err := f.VFS.Root()
	if err != nil {
		return nil, err
	}
	return &Dir{root, f}, nil
}
func (s *storageStruct) CheckAndMountAll() error {
	fmt.Println(rconfig.LoadedData().GetSectionList())
	mo := mountlib.Options{DeviceName: "a624669980_dropbox_1676528086"}
	a, e := fs.NewFs(context.TODO(), "/root/.config/rclone/rclone.conf")
	fmt.Println(e)
	aaa := func(VFS *vfs.VFS, mountpoint string, opt *mountlib.Options) (<-chan error, func() error, error) {
		f := VFS.Fs()
		fs.Debugf(f, "Mounting on %q", mountpoint)
		c, err := fuse.Mount(mountpoint, mountOptions(VFS, opt.DeviceName, opt)...)
		if err != nil {
			return nil, nil, err
		}

		filesys := NewFS(VFS, opt)
		filesys.server = fusefs.New(c, nil)

		// Serve the mount point in the background returning error to errChan
		errChan := make(chan error, 1)
		go func() {
			err := filesys.server.Serve(filesys)
			closeErr := c.Close()
			if err == nil {
				err = closeErr
			}
			errChan <- err
		}()

		// unmount := func() error {
		// 	// Shutdown the VFS
		// 	filesys.VFS.Shutdown()
		// 	return fuse.Unmount(mountpoint)
		// }
		return nil, nil, nil
	}
	mt := mountlib.NewMountPoint(aaa, "/mnt/test", a, &mo, &vfscommon.Options{})
	d, e := mt.Mount()
	fmt.Println(d, e)

	if e != nil {
		logger.Error("when CheckAndMountAll then", zap.Error(e))
		return e
	}
	defer mt.Unmount()

	storages, err := MyService.Storage().GetStorages()
	if err != nil {
		return err
	}
	logger.Info("when CheckAndMountAll storages", zap.Any("storages", storages))

	section := rconfig.LoadedData().GetSectionList()
	if err != nil {
		return err
	}
	logger.Info("when CheckAndMountAll section", zap.Any("section", section))
	for _, v := range section {
		mountPoint, found := rconfig.LoadedData().GetValue(v, "mount_point")

		if !found && len(mountPoint) == 0 {
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
			logger.Info("when CheckAndMountAll MountStorage", zap.String("mountPoint", mountPoint), zap.String("fs", v))
			err := MyService.Storage().MountStorage(mountPoint, v+":")
			if err != nil {
				logger.Error("when CheckAndMountAll then", zap.String("mountPoint", mountPoint), zap.String("fs", v), zap.Error(err))
			}
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
func (s *storageStruct) GetConfig() (httper.RemotesResult, error) {
	section, err := httper.GetAllConfigName()
	if err != nil {
		return httper.RemotesResult{}, err
	}
	return section, nil
}
func NewStorageService() StorageService {
	return &storageStruct{}
}
