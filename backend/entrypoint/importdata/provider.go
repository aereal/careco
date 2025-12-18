package importdata

import (
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/infra/firestore"
	"careco/backend/infra/gcp"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/usecases"
	"careco/backend/usecases/interactions"
	"careco/backend/usecases/ports"

	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	Provider = wire.NewSet(
		config.ProvideEnvironment,
		firestore.ClientProvider,
		firestore.ProvideDrivingRecordRepository,
		interactions.ProvideImportData,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		ProvideEntrypoint,
		providers.ProvideExportFileName,
		providers.ProvideFirestoreEmulatorAddr,
		providers.ProvideLogLevel,
		providers.ProvideServiceVersionFromGitRevision,
		wire.Bind(new(ports.DrivingRecordBulkWriter), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Bind(new(usecases.ImportData), new(*interactions.ImportData)),
		wire.Value(gcp.ProjectID("dummy")),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
)
