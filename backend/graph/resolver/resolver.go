//go:generate ./generate.bash

package resolver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"careco/backend/domain"
	"careco/backend/graph/dtos"

	"github.com/aereal/optional"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func ProvideResolver(recordCommand domain.DrivingRecordCommand, recordQuery domain.DrivingRecordQuery, tireCommand domain.TireCommand, tireQuery domain.TireQuery) *Resolver {
	return &Resolver{
		drivingRecordCommand: recordCommand,
		drivingRecordQuery:   recordQuery,
		tireCommand:          tireCommand,
		tireQuery:            tireQuery,
	}
}

type Resolver struct {
	drivingRecordCommand domain.DrivingRecordCommand
	drivingRecordQuery   domain.DrivingRecordQuery
	tireCommand          domain.TireCommand
	tireQuery            domain.TireQuery
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

func toDailyReport(record *domain.DrivingRecord) *dtos.DailyReport {
	tireID, _ := optional.Unwrap(record.TireID)
	return &dtos.DailyReport{
		OdometerValue: int(record.OdometerValue),
		RecordedAt:    record.Date,
		ReportDate:    record.Date,
		Memo:          record.Memo.Ptr(),
		TireID:        tireID,
	}
}

func toTireDTO(tire *domain.Tire) *dtos.Tire {
	return &dtos.Tire{
		ID:          string(tire.ID),
		Name:        tire.Name,
		PurchasedOn: tire.PurchasedOn.Ptr(),
		Retired:     tire.Retired,
	}
}
