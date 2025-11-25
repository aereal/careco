//go:build wireinject

package main

import (
	"context"

	"careco/backend/cmd/import-data/internal"
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/infra/firestore"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/usecases"
	"careco/backend/usecases/interactions"
	"careco/backend/usecases/ports"

	sdk "cloud.google.com/go/firestore"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func build(_ context.Context) (*internal.Entrypoint, error) {
	wire.Build(
		config.ProvideEnvironment,
		firestore.ClientProvider,
		firestore.ProvideDrivingRecordRepository,
		interactions.ProvideImportData,
		internal.ProvideEntrypoint,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		providers.ProvideExportFileName,
		providers.ProvideFirestoreEmulatorAddr,
		providers.ProvideLogLevel,
		providers.ProvideServiceVersionFromGitRevision,
		wire.Bind(new(firestore.TransactionRunner), new(*sdk.Client)),
		wire.Bind(new(ports.DrivingRecordBulkWriter), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Bind(new(usecases.ImportData), new(*interactions.ImportData)),
		wire.Value(firestore.ProjectID("dummy")),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
	return nil, nil
}
