package main

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
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

func TestMetrics(t *testing.T) {
	h := newHandler(&fakeBackend{}, "secret", time.Hour)
	get := func(path, auth string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", auth)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		return w
	}

	hits := testutil.ToFloat64(readCacheRequests.WithLabelValues("hit"))
	get("/v1/projects", "Bearer secret")
	get("/v1/projects", "Bearer secret")
	if got := testutil.ToFloat64(readCacheRequests.WithLabelValues("hit")) - hits; got != 1 {
		t.Errorf("read cache hits: got %v, want 1", got)
	}

	observeTogglCall("GetProjects", time.Millisecond, nil)
	observeTogglCall("GetProjects", time.Millisecond, errors.New("toggl down"))
	if got := testutil.ToFloat64(togglCallErrors.WithLabelValues("GetProjects")); got != 1 {
		t.Errorf("toggl call errors: got %v, want 1", got)
	}

	get("/v1/projects", "Bearer wrong")
	w := get("/metrics", "")
	if w.Code != http.StatusOK {
		t.Fatalf("/metrics without auth: status %d, want 200", w.Code)
	}
	body := w.Body.String()
	for _, want := range []string{
		`backend_toggl_http_request_duration_seconds_count{route="GET /v1/projects",status="200"}`,
		`backend_toggl_http_request_duration_seconds_count{route="unmatched",status="401"}`,
		`timenote_toggl_api_call_duration_seconds_count{method="GetProjects"} 2`,
		"go_goroutines",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("/metrics misses %q", want)
		}
	}
}
