package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithURLPath("v1/traces"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Printf("failed to create exporter: %v", err)
		return nil, err
	}

	// Creating the resource for sdktrace
	// This is where we can set service name and other attributes
	// This is optional but a recommended practice
	// This will help us identifying the service in Jaeger Backend

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	log.Println("Tracer initialized with OTLP exporter")
	return tp, nil
}
