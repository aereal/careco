package test

import (
	"careco/backend/domain/mock"

	"github.com/99designs/gqlgen/graphql/handler"
)

func provideTestServer(srv *handler.Server, drivingRecordCmd *mock.MockDrivingRecordCommand, drivingRecordQuery *mock.MockDrivingRecordQuery, tireCmd *mock.MockTireCommand, tireQuery *mock.MockTireQuery) *TestHandler {
	return &TestHandler{
		Server:               srv,
		DrivingRecordCommand: drivingRecordCmd,
		DrivingRecordQuery:   drivingRecordQuery,
		TireCommand:          tireCmd,
		TireQuery:            tireQuery,
	}
}

type TestHandler struct {
	*handler.Server

	DrivingRecordCommand *mock.MockDrivingRecordCommand
	DrivingRecordQuery   *mock.MockDrivingRecordQuery
	TireCommand          *mock.MockTireCommand
	TireQuery            *mock.MockTireQuery
}
