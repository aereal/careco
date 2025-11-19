//go:build wireinject

package test

import (
	"context"

	"careco/backend/config"
	"careco/backend/infra/firestore"

	sdk "cloud.google.com/go/firestore"
	"github.com/google/wire"
)

func BuildDrivingRecordRepository(_ context.Context) (*firestore.DrivingRecordRepository, error) {
	wire.Build(
		config.ProvideEnvironment,
		firestore.ProvideDrivingRecordRepository,
		firestore.ProvideEmulatorClient,
		provideEmulatorAddr,
		provideTestCollectionProvider,
		provideTracerProvider,
		wire.Bind(new(firestore.TransactionRunner), new(*sdk.Client)),
		wire.Value(firestore.ProjectID("test")),
	)
	return nil, nil
}
