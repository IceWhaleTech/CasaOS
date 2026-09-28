package route

import (
	"net"
	"net/http"
	"strings"
)

// isLocalRequest allows direct calls from local services and calls forwarded
// by the gateway for local clients. The gateway appends the actual client to
// X-Forwarded-For; every address must be local so a forged earlier hop cannot
// turn a remote request into an anonymous local call.
func isLocalRequest(r *http.Request) bool {
	peer, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !isLoopbackIP(peer) {
		return false
	}

	// Reject other forwarding formats rather than interpreting an untrusted
	// client-supplied address as an exemption.
	if r.Header.Get("Forwarded") != "" || r.Header.Get("X-Real-IP") != "" {
		return false
	}
	values := r.Header.Values("X-Forwarded-For")
	if len(values) == 0 && (r.Header.Get("X-Forwarded-Host") != "" || r.Header.Get("X-Forwarded-Proto") != "") {
		return false
	}
	for _, value := range values {
		for _, address := range strings.Split(value, ",") {
			if !isLoopbackIP(strings.TrimSpace(address)) {
				return false
			}
		}
	}
	return true
}

func isLoopbackIP(address string) bool {
	ip := net.ParseIP(strings.Trim(address, "[]"))
	return ip != nil && ip.IsLoopback()
}
