//go:build wireinject

package test

import (
	"testing"

	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/infra/firestore"
	"careco/backend/infra/gcp"

	"github.com/google/wire"
)

func BuildDrivingRecordRepository(_ *testing.T) (*firestore.DrivingRecordRepository, error) {
	wire.Build(
		config.ProvideEnvironment,
		firestore.ProvideClient,
		firestore.ProvideDrivingRecordRepository,
		provideContext,
		provideDatabaseID,
		providers.ProvideFirestoreEmulatorAddr,
		provideTracerProvider,
		wire.Value(gcp.ProjectID("test")),
	)
	return nil, nil
}
