package test

import (
	"careco/backend/domain/mock"

	"github.com/99designs/gqlgen/graphql/handler"
)

func provideTestServer(srv *handler.Server, drivingRecordCmd *mock.MockDrivingRecordCommand, drivingRecordQuery *mock.MockDrivingRecordQuery) *TestHandler {
	return &TestHandler{
		Server:               srv,
		DrivingRecordCommand: drivingRecordCmd,
		DrivingRecordQuery:   drivingRecordQuery,
	}
}

type TestHandler struct {
	*handler.Server

	DrivingRecordCommand *mock.MockDrivingRecordCommand
	DrivingRecordQuery   *mock.MockDrivingRecordQuery
}
