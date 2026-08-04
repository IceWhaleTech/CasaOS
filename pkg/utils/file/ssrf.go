package file

import (
	"net"
	"net/url"
	"strings"
)

// blockedHostnames contains hostnames that should never be accessed via SSRF.
var blockedHostnames = []string{
	"localhost",
	"metadata.google.internal",
}

// IsPrivateOrReservedIP checks if a hostname resolves to a private, reserved,
// or otherwise non-public IP address.
func IsPrivateOrReservedIP(hostname string) bool {
	for _, blocked := range blockedHostnames {
		if hostname == blocked {
			return true
		}
	}

	// Parse as IP address
	ip := net.ParseIP(hostname)
	if ip != nil {
		return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
	}

	// Check common private IP prefix patterns for non-standard representations
	for _, prefix := range []string{"10.", "172.", "192.168."} {
		if strings.HasPrefix(hostname, prefix) {
			return true
		}
	}
	return false
}

// IsAllowedURL validates a URL for SSRF protection.
// Returns an error if the URL points to an internal/private resource
// or uses a disallowed scheme.
func IsAllowedURL(rawURL string, allowedSchemes ...string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	hostname := parsedURL.Hostname()
	if IsPrivateOrReservedIP(hostname) {
		return errBlockedInternalResource
	}

	if len(allowedSchemes) > 0 {
		schemeAllowed := false
		for _, s := range allowedSchemes {
			if parsedURL.Scheme == s {
				schemeAllowed = true
				break
			}
		}
		if !schemeAllowed {
			return errSchemeNotAllowed
		}
	}

	return nil
}
