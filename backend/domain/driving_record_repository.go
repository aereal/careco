package domain

import (
	"context"
	"time"
)

type DrivingRecordCommand interface {
	RecordDrivingRecord(ctx context.Context, record *DrivingRecord) error
}

type DrivingRecordQuery interface {
	FindRecentRecords(ctx context.Context, first int) ([]*DrivingRecord, error)
	CalculateTotalDistance(ctx context.Context, searchPeriod Interval[time.Time]) (int64, error)
}
