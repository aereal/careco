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

	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	Provider = wire.NewSet(
		authz.ProvideAuth0Middleware,
		config.ProvideEnvironment,
		firestore.ClientProvider,
		firestore.ProvideDrivingRecordRepository,
		graph.ProvideServer,
		http.ProvideClient,
		log.ProvideGlobalInstrumentation,
		log.ProvideJSONLogger,
		log.ProvideStdoutOutput,
		o11y.ProvideResource,
		o11y.ProvideTracerProvider,
		ProvideEntrypoint,
		providers.ProvideAudience,
		providers.ProvideFirestoreEmulatorAddr,
		providers.ProvideIssuer,
		providers.ProvideLogLevel,
		providers.ProvidePort,
		providers.ProvideServiceVersionFromGitRevision,
		resolver.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(domain.DrivingRecordCommand), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(domain.DrivingRecordQuery), new(*firestore.DrivingRecordRepository)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(gcp.ProjectID("dummy")),
		wire.Value(o11y.DeploymentEnvironmentName("local")),
	)
)
