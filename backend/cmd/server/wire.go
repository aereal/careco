//go:build wireinject

package main

import (
	"context"

	"careco/backend/cmd/server/internal"
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/graph"
	"careco/backend/graph/resolver"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/web"

	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func build(_ context.Context) (*internal.Entrypoint, error) {
	wire.Build(
		config.ProvideEnvironment,
		graph.ProvideServer,
		internal.ProvideEntrypoint,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		providers.ProvideLogLevel,
		providers.ProvidePort,
		providers.ProvideServiceVersionFromGitRevision,
		resolver.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
	return nil, nil
}
