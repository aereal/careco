package domain

import "context"

type DrivingRecordCommand interface {
	RecordDrivingRecord(ctx context.Context, record *DrivingRecord) error
}

type DrivingRecordQuery interface {
	FindRecentRecords(ctx context.Context, first int) ([]*DrivingRecord, error)
}
