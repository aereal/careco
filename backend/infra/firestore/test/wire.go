//go:build wireinject

package test

import (
	"testing"

	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/infra/firestore"

	sdk "cloud.google.com/go/firestore"
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
		wire.Bind(new(firestore.TransactionRunner), new(*sdk.Client)),
		wire.Value(firestore.ProjectID("test")),
	)
	return nil, nil
}
