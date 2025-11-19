package ports

import (
	"context"

	"careco/backend/domain"
)

type DrivingRecordBulkWriter interface {
	BulkWriteDrivingRecords(ctx context.Context, items []*domain.DrivingRecord) error
}
