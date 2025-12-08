//go:build !development

package providers

import (
	"context"

	"careco/backend/config"
	"careco/backend/config/providers/live"
	"careco/backend/infra/firestore"
	"careco/backend/o11y"
)

func provideFirestoreDatabaseID(e *config.Environment) (firestore.DatabaseID, error) {
	return live.ProvideFirestoreDatabaseID(e)
}

func provideGoogleProjectID(e *config.Environment) (firestore.ProjectID, error) {
	return live.ProvideGoogleProjectID(e)
}

func provideDeploymentEnvironmentName(e *config.Environment) (o11y.DeploymentEnvironmentName, error) {
	return live.ProvideDeploymentEnvironmentName(e)
}

func provideServiceVersion(ctx context.Context, e *config.Environment) (o11y.ServiceVersion, error) {
	return live.ProvideServiceVersion(ctx, e)
}
