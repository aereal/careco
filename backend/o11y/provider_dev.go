package o11y

import (
	"context"

	"careco/backend/o11y/dev"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

func provideTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return dev.ProvideSidecarCollectorExporter(ctx)
}
