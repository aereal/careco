package live

import (
	"context"

	"careco/backend/config"
	"careco/backend/infra/firestore"
	"careco/backend/infra/gcp"
	"careco/backend/o11y"
)

func ProvideFirestoreDatabaseID(e *config.Environment) (firestore.DatabaseID, error) {
	cast := config.Cast(config.StringAs[firestore.DatabaseID])
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"FIRESTORE_DATABASE_ID",
		retrieve,
	)
}

func ProvideGoogleProjectID(e *config.Environment) (gcp.ProjectID, error) {
	cast := config.Cast(config.StringAs[gcp.ProjectID])
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"GOOGLE_CLOUD_PROJECT",
		retrieve,
	)
}

func ProvideDeploymentEnvironmentName(_ *config.Environment) (o11y.DeploymentEnvironmentName, error) {
	return "live", nil
}

func ProvideServiceVersion(_ context.Context, e *config.Environment) (o11y.ServiceVersion, error) {
	cast := config.Cast(config.StringAs[o11y.ServiceVersion])
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"OTEL_SERVICE_VERSION",
		retrieve,
	)
}
