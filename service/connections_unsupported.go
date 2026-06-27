//go:build !linux
// +build !linux

package service

import "runtime"

func (s *connectionsStruct) MountSmaba(username, host, directory, port, mountPoint, password string) error {
	return unsupportedMountError()
}

func (s *connectionsStruct) UnmountSmaba(mountPoint string) error {
	return unsupportedMountError()
}

func unsupportedMountError() error {
	return &mountUnsupportedError{platform: runtime.GOOS}
}

type mountUnsupportedError struct {
	platform string
}

func (e *mountUnsupportedError) Error() string {
	return "samba mounts are only supported on linux, not " + e.platform
}
