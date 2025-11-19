package firestore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"strconv"
	"time"

	"careco/backend/domain"
	"careco/backend/o11y/traceutils"
	"careco/backend/types"
	"careco/backend/usecases/ports"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/aereal/optional"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/iterator"
)

func ProvideDrivingRecordRepository(tp trace.TracerProvider, cp CollectionProvider, tr TransactionRunner) *DrivingRecordRepository {
	return &DrivingRecordRepository{
		collections: cp,
		tracer:      tp.Tracer("careco/backend/infra/firestore.DrivingRecordRepository"),
		txRunner:    tr,
	}
}

type DrivingRecordRepository struct {
	tracer      trace.Tracer
	collections CollectionProvider
	txRunner    TransactionRunner
}

var (
	_ domain.DrivingRecordCommand   = (*DrivingRecordRepository)(nil)
	_ domain.DrivingRecordQuery     = (*DrivingRecordRepository)(nil)
	_ ports.DrivingRecordBulkWriter = (*DrivingRecordRepository)(nil)
)

func (r *DrivingRecordRepository) RecordDrivingRecord(ctx context.Context, record *domain.DrivingRecordToRecord) (err error) {
	ctx, span := r.tracer.Start(ctx, "RecordDrivingRecord")
	defer func() { traceutils.FinishSpan(span, err) }()

	newDoc := r.collections.DrivingRecords(ctx).Doc(epochID(record.Date))
	query := r.collections.DrivingRecords(ctx).OrderBy("Date", firestore.Desc).Limit(1)
	f := func(ctx context.Context, tx *firestore.Transaction) error {
		docs := tx.Documents(query)
		defer docs.Stop()
		lastTotalDistance, err := getLastTotalDistance(docs)
		if err != nil {
			return err
		}
		newDTO := &dtoDrivingRecord{
			Date:          record.Date,
			Distance:      record.DistanceKilometers,
			TotalDistance: lastTotalDistance,
			Memo:          record.Memo.Ptr(),
		}
		if err := tx.Set(newDoc, newDTO); err != nil {
			return fmt.Errorf("firestore.Transaction.Set: %w", err)
		}
		return nil
	}
	if err := r.txRunner.RunTransaction(ctx, f); err != nil {
		return fmt.Errorf("RunTransaction: %w", err)
	}
	return nil
}

func (r *DrivingRecordRepository) CalculateTotalDistance(ctx context.Context, searchPeriod domain.Interval[time.Time]) (_ int64, err error) {
	ctx, span := r.tracer.Start(ctx, "CalculateTotalDistance")
	defer func() { traceutils.FinishSpan(span, err) }()

	totalPath := "total"
	query := r.collections.DrivingRecords(ctx).Query
	if !searchPeriod.Start.Value.IsZero() && !searchPeriod.End.Value.IsZero() {
		query = query.WhereEntity(toWhere(searchPeriod))
	}
	ret, err := query.NewAggregationQuery().WithSum("Distance", totalPath).Get(ctx)
	if err != nil {
		return 0, err
	}
	val, err := types.Cast[*firestorepb.Value](ret[totalPath])
	if err != nil {
		return 0, err
	}
	return val.GetIntegerValue(), nil
}

func (r *DrivingRecordRepository) FindRecordsInPeriod(ctx context.Context, searchPeriod domain.Interval[time.Time], direction domain.OrderDirection, limit optional.Option[int]) (_ []*domain.DrivingRecord, err error) {
	ctx, span := r.tracer.Start(ctx, "FindRecordsInPeriod")
	defer func() { traceutils.FinishSpan(span, err) }()

	query := r.collections.DrivingRecords(ctx).Query
	if whereClause := toWhere(searchPeriod); whereClause != nil {
		query = query.WhereEntity(whereClause)
	}
	if l, ok := optional.Unwrap(limit); ok {
		query = query.Limit(l)
	}
	docs := query.
		OrderBy("Date", directionMapping[direction]).
		Documents(ctx)
	records := make([]*domain.DrivingRecord, 0)
	for ret := range iterateDrivingRecords(docs) {
		record, err := ret.Value()
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (r *DrivingRecordRepository) BulkWriteDrivingRecords(ctx context.Context, items []*domain.DrivingRecord) (err error) {
	ctx, span := r.tracer.Start(ctx, "BulkWriteDrivingRecords")
	defer func() { traceutils.FinishSpan(span, err) }()

	f := func(ctx context.Context, tx *firestore.Transaction) error {
		for _, item := range items {
			data := &dtoDrivingRecord{
				Date:          item.Date,
				Distance:      item.DistanceKilometers,
				Memo:          item.Memo.Ptr(),
				TotalDistance: item.TotalDistanceKilometers,
			}
			docRef := r.collections.DrivingRecords(ctx).Doc(epochID(item.Date))
			if err := tx.Set(docRef, data); err != nil {
				return fmt.Errorf("firestore.Transaction.Set: %w", err)
			}
		}
		return nil
	}
	if err := r.txRunner.RunTransaction(ctx, f); err != nil {
		return fmt.Errorf("RunTransaction: %w", err)
	}
	return nil
}

type dtoDrivingRecord struct {
	Date          time.Time
	Distance      int64
	TotalDistance int64
	Memo          *string
}

func epochID(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

type result[T any] struct {
	value T
	err   error
}

func (r *result[T]) Value() (T, error) {
	if r.err != nil {
		return *new(T), r.err
	}
	return r.value, nil
}

func iterateDrivingRecords(docs *firestore.DocumentIterator) iter.Seq[*result[*domain.DrivingRecord]] {
	return func(yield func(*result[*domain.DrivingRecord]) bool) {
		for {
			doc, err := docs.Next()
			if errors.Is(err, iterator.Done) {
				return
			}
			if err != nil {
				_ = yield(&result[*domain.DrivingRecord]{err: err})
				return
			}
			dto := new(dtoDrivingRecord)
			if err := doc.DataTo(dto); err != nil {
				_ = yield(&result[*domain.DrivingRecord]{err: err})
				return
			}
			record := &domain.DrivingRecord{
				Date:                    dto.Date,
				DistanceKilometers:      dto.Distance,
				Memo:                    optional.FromPtr(dto.Memo),
				TotalDistanceKilometers: dto.TotalDistance,
			}
			if !yield(&result[*domain.DrivingRecord]{value: record}) {
				return
			}
		}
	}
}

func toWhere(searchPeriod domain.Interval[time.Time]) firestore.EntityFilter {
	if searchPeriod.Start.Value.IsZero() || searchPeriod.End.Value.IsZero() {
		return nil
	}
	return firestore.AndFilter{
		Filters: []firestore.EntityFilter{
			filterFragment("Date", ">", searchPeriod.Start),
			filterFragment("Date", "<", searchPeriod.End),
		},
	}
}

func filterFragment(path string, baseOp string, endpoint domain.Endpoint[time.Time]) firestore.EntityFilter {
	op := baseOp
	if !endpoint.Open {
		op += "="
	}
	return firestore.PropertyFilter{
		Path:     path,
		Operator: op,
		Value:    endpoint.Value,
	}
}

var directionMapping = map[domain.OrderDirection]firestore.Direction{
	domain.OrderDirectionAsc:  firestore.Asc,
	domain.OrderDirectionDesc: firestore.Desc,
}

func getLastTotalDistance(docs *firestore.DocumentIterator) (int64, error) {
	snapshot, err := docs.Next()
	if errors.Is(err, iterator.Done) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("firestore.DocumentIterator.Next: %w", err)
	}
	var dto struct{ TotalDistance int64 }
	if err := snapshot.DataTo(&dto); err != nil {
		return 0, fmt.Errorf("firestore.DocumentSnapshot.DataTo: %w", err)
	}
	return dto.TotalDistance, nil
}
