//go:generate ./generate.bash

package resolver

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *Resolver) getOdometerValue(ctx context.Context, period domain.Interval[time.Time]) (int64, error) {
	record, err := r.drivingRecordQuery.FindLastRecordInPeriod(ctx, period)
	if errors.Is(err, domain.ErrDrivingRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("FindLastRecordInPeriod: %w", err)
	}
	return record.OdometerValue, nil
}
