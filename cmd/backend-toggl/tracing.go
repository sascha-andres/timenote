package main

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// setupTracing installs the global tracer provider and propagator. The
// standard OTEL_* environment variables configure them (exporter, endpoint,
// sampler, resource, propagators). Tracing is off if OTEL_TRACES_EXPORTER is
// not set. The returned function flushes and stops the tracer provider.
func setupTracing(ctx context.Context) (func(context.Context) error, error) {
	// the spec default is "otlp"; off by default prevents export errors when no collector runs
	if os.Getenv("OTEL_TRACES_EXPORTER") == "" {
		return func(context.Context) error { return nil }, nil
	}
	exporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())
	return tp.Shutdown, nil
}
