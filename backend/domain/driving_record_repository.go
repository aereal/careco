package domain

import (
	"context"
	"time"

	"github.com/aereal/optional"
)

type DrivingRecordCommand interface {
	RecordDrivingRecord(ctx context.Context, record *DrivingRecord) error
}

type DrivingRecordQuery interface {
	FindRecordsInPeriod(ctx context.Context, searchPeriod Interval[time.Time], direction OrderDirection, limit optional.Option[int]) ([]*DrivingRecord, error)
	FindLastRecordInPeriod(ctx context.Context, searchPeriod Interval[time.Time]) (*DrivingRecord, error)
}
