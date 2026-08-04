package file

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrPathTraversal          = errors.New("path traversal detected")
	ErrEmptyPath              = errors.New("path is empty")
	ErrPathNotAllowed         = errors.New("path is not within allowed directories")
	errBlockedInternalResource = errors.New("access to internal resources is forbidden")
	errSchemeNotAllowed       = errors.New("only http and https schemes are allowed")
)

// AllowedBaseDirs defines the directories that file operations are confined to.
// This prevents path traversal attacks. Add any additional base directories as needed.
var AllowedBaseDirs = []string{
	"/DATA",
	"/var/lib/casaos",
	"/tmp",
	"/etc/samba",
}

// SanitizePath validates the given path and ensures it stays within one of the
// AllowedBaseDirs. It rejects any path containing ".." components.
func SanitizePath(rawPath string) (string, error) {
	if len(rawPath) == 0 {
		return "", ErrEmptyPath
	}

	// Reject obvious traversal attempts early
	if strings.Contains(rawPath, "..") {
		return "", ErrPathTraversal
	}

	// Clean the path to normalize it
	cleaned := filepath.Clean(rawPath)

	// Resolve symlinks to prevent symlink-based traversal
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}

	// Check against allowed base directories
	for _, base := range AllowedBaseDirs {
		absBase, err := filepath.Abs(base)
		if err != nil {
			continue
		}
		// Ensure the base directory ends with separator for accurate prefix matching
		if !strings.HasSuffix(absBase, string(filepath.Separator)) {
			absBase += string(filepath.Separator)
		}
		if strings.HasPrefix(absPath, absBase) || absPath == strings.TrimSuffix(absBase, string(filepath.Separator)) {
			return absPath, nil
		}
	}

	return "", ErrPathNotAllowed
}

// SanitizePaths validates a list of paths, returning only valid ones and the first error encountered.
func SanitizePaths(paths []string) ([]string, error) {
	sanitized := make([]string, 0, len(paths))
	for _, p := range paths {
		s, err := SanitizePath(p)
		if err != nil {
			return nil, err
		}
		sanitized = append(sanitized, s)
	}
	return sanitized, nil
}
