package file

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrPathTraversal           = errors.New("path traversal detected")
	ErrEmptyPath               = errors.New("path is empty")
	ErrPathNotAllowed          = errors.New("path is not within allowed directories")
	errBlockedInternalResource = errors.New("access to internal resources is forbidden")
	errSchemeNotAllowed        = errors.New("only http and https schemes are allowed")
	ErrTrashNotAllowed         = errors.New("trash directory is not within allowed directories")
)

// AllowedBaseDirs defines the directories that file operations are confined to.
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
	if strings.Contains(rawPath, "..") {
		return "", ErrPathTraversal
	}
	cleaned := filepath.Clean(rawPath)
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}
	for _, base := range AllowedBaseDirs {
		absBase, err := filepath.Abs(base)
		if err != nil {
			continue
		}
		if !strings.HasSuffix(absBase, string(filepath.Separator)) {
			absBase += string(filepath.Separator)
		}
		if strings.HasPrefix(absPath, absBase) || absPath == strings.TrimSuffix(absBase, string(filepath.Separator)) {
			return absPath, nil
		}
	}
	return "", ErrPathNotAllowed
}

// SanitizePaths validates a list of paths.
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

// SoftDelete moves a file or directory to the trash directory instead of permanently deleting it.
// The trash is organized as .casaos-trash/<timestamp>_<basename> to avoid collisions.
func SoftDelete(sourcePath string) error {
	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}

	// Determine the base directory containing the source
	baseDir := filepath.Dir(absSource)
	baseName := filepath.Base(absSource)

	// Create trash directory within the same filesystem
	trashDir := filepath.Join(baseDir, ".casaos-trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create trash directory: %w", err)
	}

	// Verify trash directory is within allowed dirs
	trashDir, err = SanitizePath(trashDir)
	if err != nil {
		return ErrTrashNotAllowed
	}

	// Generate unique trash name with timestamp
	trashName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), baseName)
	trashPath := filepath.Join(trashDir, trashName)

	// Move to trash (atomic rename if on same filesystem)
	if err := os.Rename(absSource, trashPath); err != nil {
		// If rename fails (different filesystem), fall back to copy+delete
		if err := copyToTrash(absSource, trashPath); err != nil {
			return fmt.Errorf("failed to move to trash: %w", err)
		}
		// Remove original after successful copy
		if err := os.RemoveAll(absSource); err != nil {
			return fmt.Errorf("failed to remove original after trash copy: %w", err)
		}
	}

	return nil
}

// copyToTrash copies a file or directory to the trash path.
func copyToTrash(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = out.ReadFrom(in); err != nil {
		return err
	}

	// Preserve file permissions
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

// IsPrivateOrReservedIP checks if a hostname resolves to a private, reserved,
// or otherwise non-public IP address.
func IsPrivateOrReservedIP(hostname string) bool {
	for _, blocked := range []string{"localhost", "metadata.google.internal"} {
		if hostname == blocked {
			return true
		}
	}
	ip := net.ParseIP(hostname)
	if ip != nil {
		return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
	}
	return false
}

// IsAllowedURL validates a URL for SSRF protection.
func IsAllowedURL(rawURL string, allowedSchemes ...string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if IsPrivateOrReservedIP(parsedURL.Hostname()) {
		return errBlockedInternalResource
	}
	for _, s := range allowedSchemes {
		if parsedURL.Scheme == s {
			return nil
		}
	}
	if len(allowedSchemes) > 0 {
		return errSchemeNotAllowed
	}
	return nil
}
