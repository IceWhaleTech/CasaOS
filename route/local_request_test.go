package route

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsLocalRequest(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		want       bool
	}{
		{"direct local", "127.0.0.1:12345", nil, true},
		{"direct IPv6 local", "[::1]:12345", nil, true},
		{"remote claiming local", "192.0.2.10:12345", map[string]string{"X-Forwarded-For": "127.0.0.1"}, false},
		{"gateway external client", "127.0.0.1:12345", map[string]string{"X-Forwarded-For": "192.0.2.10"}, false},
		{"gateway local client", "127.0.0.1:12345", map[string]string{"X-Forwarded-For": "127.0.0.1"}, true},
		{"injected first hop", "127.0.0.1:12345", map[string]string{"X-Forwarded-For": "127.0.0.1, 192.0.2.10"}, false},
		{"conflicting real IP", "127.0.0.1:12345", map[string]string{"X-Forwarded-For": "127.0.0.1", "X-Real-IP": "192.0.2.10"}, false},
		{"untrusted forwarded header", "127.0.0.1:12345", map[string]string{"Forwarded": "for=127.0.0.1"}, false},
		{"proxy without client address", "127.0.0.1:12345", map[string]string{"X-Forwarded-Proto": "http"}, false},
		{"malformed forwarded address", "127.0.0.1:12345", map[string]string{"X-Forwarded-For": "unknown"}, false},
		{"malformed peer address", "127.0.0.1", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/v1/sys/hardware", nil)
			r.RemoteAddr = tt.remoteAddr
			for name, value := range tt.headers {
				r.Header.Set(name, value)
			}
			if got := isLocalRequest(r); got != tt.want {
				t.Errorf("isLocalRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
