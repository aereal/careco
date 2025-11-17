package ports

import (
	"context"

	"careco/backend/domain"
)

type DrivingRecordBulkWriter interface {
	BulkWriteDrivingRecords(ctx context.Context, records []*domain.DrivingRecord) error
}
