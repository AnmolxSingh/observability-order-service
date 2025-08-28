package otel

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// InitLogger sets up a structured logger that sends logs to the OTEL Collector.
func InitLogger(res *resource.Resource) (func(context.Context) error, error) {
	// Create the OTLP log exporter
	logExporter, err := otlploghttp.New(context.Background(),
		otlploghttp.WithEndpoint("localhost:4318"),
		otlploghttp.WithURLPath("v1/logs"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create a LoggerProvider
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)

	// Set the global logger provider
	global.SetLoggerProvider(loggerProvider)

	// Use a custom handler to automatically add trace context to logs
	handler := NewOtelSlogHandler(slog.NewJSONHandler(os.Stdout, nil))
	logger := slog.New(handler)

	// Set the new logger as the global default
	slog.SetDefault(logger)

	return loggerProvider.Shutdown, nil
}

// OtelSlogHandler is a custom slog.Handler that adds trace context to logs.
type OtelSlogHandler struct {
	slog.Handler
}

func NewOtelSlogHandler(handler slog.Handler) *OtelSlogHandler {
	return &OtelSlogHandler{Handler: handler}
}

// Handle adds trace_id and span_id to the log record if a span is active.
func (h *OtelSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}
