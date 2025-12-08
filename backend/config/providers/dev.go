//go:build development

package providers

import (
	"context"

	"careco/backend/config"
	"careco/backend/config/providers/dev"
	"careco/backend/infra/firestore"
	"careco/backend/o11y"
)

func provideFirestoreDatabaseID(e *config.Environment) (firestore.DatabaseID, error) {
	return dev.ProvideFirestoreDatabaseID(e)
}

func provideGoogleProjectID(e *config.Environment) (firestore.ProjectID, error) {
	return dev.ProvideGoogleProjectID(e)
}

func provideDeploymentEnvironmentName(e *config.Environment) (o11y.DeploymentEnvironmentName, error) {
	return dev.ProvideDeploymentEnvironmentName(e)
}

func provideServiceVersion(ctx context.Context, e *config.Environment) (o11y.ServiceVersion, error) {
	return dev.ProvideServiceVersion(ctx, e)
}
