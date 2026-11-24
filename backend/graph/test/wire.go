//go:build wireinject

package test

import (
	"careco/backend/domain"
	"careco/backend/domain/mock"
	"careco/backend/graph"
	"careco/backend/graph/resolver"
	"careco/backend/o11y/test"

	"github.com/google/wire"
	"go.uber.org/mock/gomock"
)

func BuildHandler(_ *gomock.Controller) *TestHandler {
	wire.Build(
		graph.ProvideServer,
		mock.NewMockDrivingRecordCommand,
		mock.NewMockDrivingRecordQuery,
		mock.NewMockTireCommand,
		mock.NewMockTireQuery,
		provideTestServer,
		resolver.ProvideResolver,
		test.ProvideNoopTracerProvider,
		wire.Bind(new(domain.DrivingRecordCommand), new(*mock.MockDrivingRecordCommand)),
		wire.Bind(new(domain.DrivingRecordQuery), new(*mock.MockDrivingRecordQuery)),
		wire.Bind(new(domain.TireCommand), new(*mock.MockTireCommand)),
		wire.Bind(new(domain.TireQuery), new(*mock.MockTireQuery)),
	)
	return nil
}
