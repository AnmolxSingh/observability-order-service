// otel/metrics.go
package otel

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// InitMetrics sets up the Prometheus metrics exporter.
// It returns the HTTP handler for the /metrics endpoint and a shutdown function.
func InitMetrics(res *resource.Resource) (http.Handler, func(context.Context) error, error) {
	// The exporter is the component that sends data to Prometheus.
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	// The MeterProvider is the entry point for the metrics SDK.
	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	// Set the global meter provider, so you can get a Meter from it globally.
	otel.SetMeterProvider(mp)
	slog.Info("Prometheus metrics exporter initialized")

	// Return the Prometheus HTTP handler and the meter provider's shutdown function.
	return promhttp.Handler(), mp.Shutdown, nil
}
