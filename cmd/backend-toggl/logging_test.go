package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogRequests(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)

	h := newHandler(&fakeBackend{}, "secret", 0)
	tests := []struct {
		auth, method, path, want string
	}{
		{"Bearer wrong", http.MethodGet, "/v1/clients", "method=GET path=/v1/clients status=401"},
		{"Bearer secret", http.MethodGet, "/v1/clients", "method=GET path=/v1/clients status=200"},
		{"Bearer secret", http.MethodPost, "/v1/cache", "method=POST path=/v1/cache status=204 bytes=0"},
	}
	for _, tt := range tests {
		buf.Reset()
		req := httptest.NewRequest(tt.method, tt.path, nil)
		req.Header.Set("Authorization", tt.auth)
		h.ServeHTTP(httptest.NewRecorder(), req)
		if got := buf.String(); !strings.Contains(got, tt.want) || !strings.Contains(got, "duration=") {
			t.Errorf("%s %s: log %q, want %q and duration", tt.method, tt.path, got, tt.want)
		}
	}
}
