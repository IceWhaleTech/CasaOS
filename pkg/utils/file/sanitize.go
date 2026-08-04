package file

import (
	"errors"
	"net"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	ErrPathTraversal           = errors.New("path traversal detected")
	ErrEmptyPath               = errors.New("path is empty")
	ErrPathNotAllowed          = errors.New("path is not within allowed directories")
	errBlockedInternalResource = errors.New("access to internal resources is forbidden")
	errSchemeNotAllowed        = errors.New("only http and https schemes are allowed")
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
