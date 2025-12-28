package server

import (
	"careco/backend/authz"
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/domain"
	"careco/backend/graph"
	"careco/backend/graph/resolver"
	"careco/backend/infra/firestore"
	"careco/backend/infra/gcp"
	"careco/backend/infra/http"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/web"

	firestoresdk "cloud.google.com/go/firestore"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	Provider = wire.NewSet(
		authz.ProvideMiddleware,
		BindGCPProject,
		config.ProvideEnvironment,
		firestore.ProvideDrivingRecordRepository,
		graph.ProvideServer,
		http.ProvideClient,
		log.ProvideGlobalInstrumentation,
		log.ProvideHandler,
		log.ProvideStdoutOutput,
		o11y.ProvideGlobalLoggerInstrumentation,
		o11y.ProvideTracerProvider,
		ProvideEntrypoint,
		providers.ProvideAudience,
		providers.ProvideFirestoreEmulatorAddr,
		providers.ProvideIssuer,
		providers.ProvideLogLevel,
		providers.ProvidePort,
		resolver.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(domain.DrivingRecordCommand), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(domain.DrivingRecordQuery), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
	)
	DevProvider = wire.NewSet(
		firestore.ProvideEmulatorClient,
		o11y.ProvideResource,
		o11y.ProvideSidecarCollectorExporter,
		Provider,
		providers.ProvideServiceVersionFromGitRevision,
		wire.Value(firestore.DatabaseID(firestoresdk.DefaultDatabaseID)),
		wire.Value(gcp.ProjectID("dummy")),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
	ProductionProvider = wire.NewSet(
		firestore.ProvideClient,
		o11y.ProvideGoogleCloudRunResource,
		o11y.ProvideGoogleTelemetryTraceExporter,
		Provider,
		providers.ProvideFirestoreDatabaseID,
		providers.ProvideGoogleProjectID,
		providers.ProvideServiceVersion,
		wire.Value(o11y.DeploymentEnvironmentName("production")),
	)
)
