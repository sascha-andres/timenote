package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "backend_toggl_http_request_duration_seconds",
		Help: "Duration of HTTP requests by route pattern and status.",
	}, []string{"route", "status"})

	readCacheRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backend_toggl_read_cache_requests_total",
		Help: "GET requests served from the read cache (hit) or from the backend (miss).",
	}, []string{"result"})

	togglCallDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "timenote_toggl_api_call_duration_seconds",
		Help: "Duration of toggl API calls by method.",
	}, []string{"method"})

	togglCallErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "timenote_toggl_api_call_errors_total",
		Help: "Failed toggl API calls by method.",
	}, []string{"method"})
)

// observeTogglCall records a toggl API call; main sets it as persistence.OnTogglCall.
func observeTogglCall(method string, duration time.Duration, err error) {
	togglCallDuration.WithLabelValues(method).Observe(duration.Seconds())
	if err != nil {
		togglCallErrors.WithLabelValues(method).Inc()
	}
}

// statusRecorder keeps the status code and the number of body bytes that a handler writes.
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// Unwrap lets http.ResponseController reach the original writer.
func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// instrument logs each request at debug level and records its duration
// after next has handled it.
func instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		if rec.status == 0 {
			rec.status = http.StatusOK
		}
		// the muxes set r.Pattern; "/" is the catch-all to the API mux, so the request matched no route
		route := r.Pattern
		if route == "" || route == "/" {
			route = "unmatched"
		}
		requestDuration.WithLabelValues(route, strconv.Itoa(rec.status)).Observe(duration.Seconds())
		slog.Debug("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"bytes", rec.bytes,
			"duration", duration,
			"remote", r.RemoteAddr,
		)
	})
}
