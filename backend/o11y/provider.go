package o11y

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type (
	ServiceVersion            string
	DeploymentEnvironmentName string
)

func ProvideResource(ctx context.Context, svcVersion ServiceVersion, depEnv DeploymentEnvironmentName) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName("careco-backend"),
			semconv.ServiceVersion(string(svcVersion)),
			semconv.DeploymentEnvironmentName(string(depEnv)),
		),
	)
}

func ProvideTracerProvider(ctx context.Context, exporter *otlptrace.Exporter, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	), nil
}

func ProvideSidecarCollectorExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}
