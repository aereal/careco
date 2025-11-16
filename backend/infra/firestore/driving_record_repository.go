package firestore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"strconv"
	"time"

	"careco/backend/domain"
	"careco/backend/types"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/aereal/optional"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/iterator"
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

var (
	_ domain.DrivingRecordCommand = (*DrivingRecordRepository)(nil)
	_ domain.DrivingRecordQuery   = (*DrivingRecordRepository)(nil)
)

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

func (r *DrivingRecordRepository) FindRecentRecords(ctx context.Context, first int) (_ []*domain.DrivingRecord, err error) {
	ctx, span := r.tracer.Start(ctx, "FindRecentRecords")
	defer span.End()

	docs := r.collections.DrivingRecords(ctx).Limit(first).OrderBy("Date", firestore.Desc).Documents(ctx)
	records := make([]*domain.DrivingRecord, 0, first)
	for ret := range iterateDrivingRecords(docs) {
		record, err := ret.Value()
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (r *DrivingRecordRepository) CalculateTotalDistance(ctx context.Context, searchPeriod domain.Interval[time.Time]) (_ int64, err error) {
	ctx, span := r.tracer.Start(ctx, "CalculateTotalDistance")
	defer span.End()

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
	defer span.End()

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

type dtoDrivingRecord struct {
	Date     time.Time
	Distance int64
	Memo     *string
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
				Date:               dto.Date,
				DistanceKilometers: dto.Distance,
				Memo:               optional.FromPtr(dto.Memo),
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
	if endpoint.Open {
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
