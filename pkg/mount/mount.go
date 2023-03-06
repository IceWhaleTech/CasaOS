package mount

import (
	"fmt"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"github.com/rclone/rclone/cmd/mountlib"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/log"
	"github.com/rclone/rclone/vfs"
)

func MountFn(VFS *vfs.VFS, mountpoint string, opt *mountlib.Options) (<-chan error, func() error, error) {

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

	unmount := func() error {
		// Shutdown the VFS
		filesys.VFS.Shutdown()
		return fuse.Unmount(mountpoint)
	}
	return errChan, unmount, nil

}
func NewFS(VFS *vfs.VFS, opt *mountlib.Options) *FS {
	fsys := &FS{
		VFS: VFS,
		f:   VFS.Fs(),
		opt: opt,
	}
	return fsys
}

type FS struct {
	*vfs.VFS
	f      fs.Fs
	opt    *mountlib.Options
	server *fusefs.Server
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

// Root returns the root node
func (f *FS) Root() (node fusefs.Node, err error) {
	defer log.Trace("", "")("node=%+v, err=%v", &node, &err)
	root, err := f.VFS.Root()
	if err != nil {
		return nil, err
	}
	return &Dir{root, f}, nil
}
