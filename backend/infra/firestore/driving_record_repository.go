package firestore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"careco/backend/domain"

	"go.opentelemetry.io/otel/trace"
)

func ProvideDrivingRecordRepository(tp trace.TracerProvider, cp CollectionProvider) *DrivingRecordRepository {
	return &DrivingRecordRepository{
		collections: cp,
		tracer:      tp.Tracer("careco/backend/infra/firestore.DrivingRecordRepository"),
	}
}

type DrivingRecordRepository struct {
	tracer      trace.Tracer
	collections CollectionProvider
}

var _ domain.DrivingRecordCommand = (*DrivingRecordRepository)(nil)

func (r *DrivingRecordRepository) RecordDrivingRecord(ctx context.Context, record *domain.DrivingRecord) (err error) {
	ctx, span := r.tracer.Start(ctx, "RecordDrivingRecord")
	defer span.End()

	utc := record.Date.UTC()
	utcDate := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	data := &dtoDrivingRecord{
		Date:     utcDate,
		Distance: record.DistanceKilometers,
		Memo:     record.Memo.Ptr(),
	}
	if _, err := r.collections.DrivingRecords(ctx).Doc(epochID(utcDate)).Set(ctx, data); err != nil {
		return fmt.Errorf("firestore.DocumentRef.Set: %w", err)
	}
	return nil
}

type dtoDrivingRecord struct {
	Date     time.Time
	Distance int64
	Memo     *string
}

func epochID(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
