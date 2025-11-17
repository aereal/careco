//go:build wireinject

package test

import (
	"context"

	"careco/backend/config"
	"careco/backend/infra/firestore"
	o11ytest "careco/backend/o11y/test"

	sdk "cloud.google.com/go/firestore"
	"github.com/google/wire"
)

func BuildDrivingRecordRepository(_ context.Context) (*firestore.DrivingRecordRepository, error) {
	wire.Build(
		config.ProvideEnvironment,
		firestore.ProvideDrivingRecordRepository,
		firestore.ProvideEmulatorClient,
		o11ytest.ProvideNoopTracerProvider,
		provideEmulatorAddr,
		provideTestCollectionProvider,
		wire.Bind(new(firestore.TransactionRunner), new(*sdk.Client)),
		wire.Value(firestore.ProjectID("test")),
	)
	return nil, nil
}
