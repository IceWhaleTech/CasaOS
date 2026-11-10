package route

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDebugRequiresAuthentication(t *testing.T) {
	handler := InitV1Router()
	for _, tc := range []struct {
		name string
		xff  string
	}{
		{name: "without forwarded header"},
		{name: "spoofed local header", xff: "127.0.0.1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/sys/debug", nil)
			request.RemoteAddr = "192.0.2.10:12345"
			if tc.xff != "" {
				request.Header.Set("X-Forwarded-For", tc.xff)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
			}
		})
	}
}
