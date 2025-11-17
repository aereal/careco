//go:generate ./generate.bash

package resolver

import (
	"careco/backend/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func ProvideResolver(recordCommand domain.DrivingRecordCommand, recordQuery domain.DrivingRecordQuery) *Resolver {
	return &Resolver{
		drivingRecordCommand: recordCommand,
		drivingRecordQuery:   recordQuery,
	}
}

type Resolver struct {
	drivingRecordCommand domain.DrivingRecordCommand
	drivingRecordQuery   domain.DrivingRecordQuery
}
