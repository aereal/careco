//go:build !development

package o11y

import (
	"context"

	"careco/backend/o11y/live"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

func provideTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return live.ProvideGoogleTelemetryTraceExporter(ctx)
}
