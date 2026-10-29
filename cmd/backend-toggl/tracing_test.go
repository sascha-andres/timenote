package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracing(t *testing.T) {
	spans := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans)))
	otel.SetTextMapPropagator(propagation.TraceContext{})

	h := newHandler(&fakeBackend{}, "secret", 0)
	for _, path := range []string{"/v1/projects", "/metrics"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer secret")
		req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		h.ServeHTTP(httptest.NewRecorder(), req)
	}

	ended := spans.Ended()
	if len(ended) != 1 {
		t.Fatalf("got %d spans, want 1 (no span for /metrics)", len(ended))
	}
	if got := ended[0].Name(); got != "GET /v1/projects" {
		t.Errorf("span name %q, want %q", got, "GET /v1/projects")
	}
	if got := ended[0].Parent().TraceID().String(); got != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("parent trace id %q, want the one from traceparent", got)
	}
}
