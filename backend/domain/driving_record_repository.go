package domain

import (
	"context"
	"time"

	"github.com/aereal/optional"
)

type DrivingRecordCommand interface {
	RecordDrivingRecord(ctx context.Context, record *DrivingRecordToRecord) error
}

type DrivingRecordQuery interface {
	CalculateTotalDistance(ctx context.Context, searchPeriod Interval[time.Time]) (int64, error)
	FindRecordsInPeriod(ctx context.Context, searchPeriod Interval[time.Time], direction OrderDirection, limit optional.Option[int]) ([]*DrivingRecord, error)
}
