package o11y

import (
	"context"
	"fmt"
	"iter"
	"slices"

	infragcp "careco/backend/infra/gcp"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

type (
	ServiceVersion            string
	DeploymentEnvironmentName string
)

func ProvideResource(ctx context.Context, svcVersion ServiceVersion, depEnv DeploymentEnvironmentName) (*resource.Resource, error) {
	return resource.New(ctx, slices.Collect(commonResourceOptions(svcVersion, depEnv))...)
}

func ProvideGoogleCloudRunResource(ctx context.Context, gcpProject infragcp.ProjectID, svcVersion ServiceVersion, depEnv DeploymentEnvironmentName) (*resource.Resource, error) {
	opts := []resource.Option{
		resource.WithAttributes(attribute.String("gcp.project_id", string(gcpProject))), // this is mandatory for Google Cloud Telemetry API
		resource.WithDetectors(gcp.NewDetector()),
	}
	opts = slices.AppendSeq(opts, commonResourceOptions(svcVersion, depEnv))
	return resource.New(ctx, opts...)
}

func commonResourceOptions(svcVersion ServiceVersion, depEnv DeploymentEnvironmentName) iter.Seq[resource.Option] {
	return func(yield func(resource.Option) bool) {
		if !yield(resource.WithTelemetrySDK()) {
			return
		}
		if !yield(resource.WithAttributes(
			semconv.ServiceName("careco-backend"),
			semconv.ServiceVersion(string(svcVersion)),
			semconv.DeploymentEnvironmentName(string(depEnv)),
		)) {
			return
		}
	}
}

func ProvideTracerProvider(ctx context.Context, exporter *otlptrace.Exporter, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	), nil
}

func ProvideSidecarCollectorExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}

func ProvideGoogleTelemetryTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	creds, err := oauth.NewApplicationDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("oauth.NewApplicationDefault: %w", err)
	}
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("telemetry.googleapis.com:443"),
		otlptracegrpc.WithDialOption(grpc.WithPerRPCCredentials(creds)),
	)
}
