package domain

import "context"

type DrivingRecordCommand interface {
	RecordDrivingRecord(ctx context.Context, record *DrivingRecord) error
}
