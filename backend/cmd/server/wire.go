//go:build wireinject

package main

import (
	"context"

	"careco/backend/cmd/server/internal"
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/domain"
	"careco/backend/graph"
	"careco/backend/graph/resolver"
	"careco/backend/infra/firestore"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/web"

	sdk "cloud.google.com/go/firestore"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func build(_ context.Context) (*internal.Entrypoint, error) {
	wire.Build(
		config.ProvideEnvironment,
		firestore.ProvideDrivingRecordRepository,
		firestore.ProvideEmulatorClient,
		firestore.ProvideProductionCollectionProvider,
		graph.ProvideServer,
		internal.ProvideEntrypoint,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		providers.ProvideFirestoreEmulatorAddr,
		providers.ProvideLogLevel,
		providers.ProvidePort,
		providers.ProvideServiceVersionFromGitRevision,
		resolver.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(domain.DrivingRecordCommand), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(domain.DrivingRecordQuery), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(firestore.TransactionRunner), new(*sdk.Client)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(firestore.ProjectID("dummy")),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
	return nil, nil
}
