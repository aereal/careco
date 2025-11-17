//go:build wireinject

package main

import (
	"context"

	"careco/backend/cmd/import-data/internal"
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/usecases"
	"careco/backend/usecases/interactions"

	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func build(_ context.Context) (*internal.Entrypoint, error) {
	wire.Build(
		config.ProvideEnvironment,
		interactions.ProvideImportData,
		internal.ProvideEntrypoint,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		providers.ProvideLogLevel,
		providers.ProvideServiceVersionFromGitRevision,
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Bind(new(usecases.ImportData), new(*interactions.ImportData)),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
	return nil, nil
}
