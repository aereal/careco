package live

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func ProvideGoogleTelemetryTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	creds, err := oauth.NewApplicationDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("oauth.NewApplicationDefault: %w", err)
	}
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("https://telemetry.googleapis.com"),
		otlptracegrpc.WithDialOption(grpc.WithPerRPCCredentials(creds)),
	)
}
