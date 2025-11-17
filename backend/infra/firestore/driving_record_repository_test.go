package firestore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"careco/backend/domain"
	"careco/backend/infra/firestore"
	"careco/backend/infra/firestore/test"
	"careco/backend/tests"

	"github.com/aereal/optional"
)

func TestDrivingRecordRepository_RecordDrivingRecord(t *testing.T) {
	t.Parallel()

	t1 := time.Now().Truncate(time.Microsecond)
	t2 := t1.Add(time.Second * -1)
	t3 := t2.Add(time.Hour * 24 * -1)
	t4 := t3.Add(time.Second * -1)
	t5 := t4.Add(time.Second * -1)
	testCases := []struct {
		name           string
		inputs         []*domain.DrivingRecord
		intervalToFind domain.Interval[time.Time]
		wantRecords    []*domain.DrivingRecord
		wantErr        error
	}{
		{
			name: "ok",
			inputs: []*domain.DrivingRecord{
				{
					Date:               t1,
					DistanceKilometers: 123,
				},
			},
			wantRecords: []*domain.DrivingRecord{
				{
					Date:               t1,
					DistanceKilometers: 123,
				},
			},
			wantErr: nil,
			intervalToFind: domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(t2),
				End:   domain.ClosedEndpoint(t1),
			},
		},
		{
			name: "multiple calls on same date",
			inputs: []*domain.DrivingRecord{
				{Date: t3, DistanceKilometers: 45},
				{Date: t4, DistanceKilometers: 67},
			},
			wantRecords: []*domain.DrivingRecord{
				{Date: t4, DistanceKilometers: 67},
				{Date: t3, DistanceKilometers: 45},
			},
			wantErr: nil,
			intervalToFind: domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(t5),
				End:   domain.ClosedEndpoint(t3),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := test.Context(t)
			r, err := test.BuildDrivingRecordRepository(ctx)
			if err != nil {
				t.Fatal(err)
			}
			var gotErr error
			for _, input := range tc.inputs {
				err := r.RecordDrivingRecord(ctx, input)
				if err != nil {
					gotErr = errors.Join(gotErr, err)
				}
			}
			tests.AssertsErrors(t, tc.wantErr, gotErr)

			gotRecords, err := r.FindRecordsInPeriod(ctx, tc.intervalToFind, domain.OrderDirectionAsc, optional.None[int]())
			if err != nil {
				t.Fatal(err)
			}
			if err := tests.Diff(tc.wantRecords, gotRecords, tests.EquateOptional[string](), tests.EquateOptional[int]()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDrivingRecordRepository_CalculateTotalDistance(t *testing.T) {
	t.Parallel()

	baseDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		interval domain.Interval[time.Time]
		wantVal  int64
		wantErr  error
		prepare  func(ctx context.Context, r *firestore.DrivingRecordRepository) error
	}{
		{
			name:     "ok/no filter",
			interval: domain.EmptyInterval[time.Time](),
			wantVal:  6,
			wantErr:  nil,
			prepare: func(ctx context.Context, r *firestore.DrivingRecordRepository) error {
				return errors.Join(
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 1, Date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 2, Date: time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 3, Date: time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC)}),
				)
			},
		},
		{
			name: "ok/with month filter",
			interval: domain.Interval[time.Time]{
				Start: domain.OpenEndpoint(baseDate),
				End:   domain.ClosedEndpoint(nextMonth),
			},
			wantVal: 5,
			wantErr: nil,
			prepare: func(ctx context.Context, r *firestore.DrivingRecordRepository) error {
				return errors.Join(
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 1, Date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 2, Date: time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 3, Date: time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC)}),
				)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := test.Context(t)
			r, err := test.BuildDrivingRecordRepository(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if tc.prepare != nil {
				if err := tc.prepare(ctx, r); err != nil {
					t.Fatal(err)
				}
			}
			got, gotErr := r.CalculateTotalDistance(ctx, tc.interval)
			tests.AssertsErrors(t, tc.wantErr, gotErr)
			if gotErr != nil {
				return
			}
			if got != tc.wantVal {
				t.Errorf("want=%d got=%d", tc.wantVal, got)
			}
		})
	}
}

func TestDrivingRecordRepository_FindRecordsInPeriod(t *testing.T) {
	t.Parallel()

	baseDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name      string
		interval  domain.Interval[time.Time]
		direction domain.OrderDirection
		limit     optional.Option[int]
		want      []*domain.DrivingRecord
		wantErr   error
		prepare   func(ctx context.Context, r *firestore.DrivingRecordRepository) error
	}{
		{
			name:     "ok/no filter",
			interval: domain.EmptyInterval[time.Time](),
			want: []*domain.DrivingRecord{
				{
					DistanceKilometers: 1,
					Date:               time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					DistanceKilometers: 3,
					Date:               time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					DistanceKilometers: 2,
					Date:               time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
			prepare: func(ctx context.Context, r *firestore.DrivingRecordRepository) error {
				return errors.Join(
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 1, Date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 2, Date: time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 3, Date: time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC)}),
				)
			},
		},
		{
			name: "ok/with month filter",
			interval: domain.Interval[time.Time]{
				Start: domain.OpenEndpoint(baseDate),
				End:   domain.ClosedEndpoint(nextMonth),
			},
			want: []*domain.DrivingRecord{
				{
					DistanceKilometers: 3,
					Date:               time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					DistanceKilometers: 2,
					Date:               time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
			prepare: func(ctx context.Context, r *firestore.DrivingRecordRepository) error {
				return errors.Join(
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 1, Date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 2, Date: time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 3, Date: time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC)}),
				)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := test.Context(t)
			r, err := test.BuildDrivingRecordRepository(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if tc.prepare != nil {
				if err := tc.prepare(ctx, r); err != nil {
					t.Fatal(err)
				}
			}
			got, gotErr := r.FindRecordsInPeriod(ctx, tc.interval, tc.direction, tc.limit)
			tests.AssertsErrors(t, tc.wantErr, gotErr)
			if gotErr != nil {
				return
			}
			if err := tests.Diff(tc.want, got, tests.EquateOptional[string]()); err != nil {
				t.Error(err)
			}
		})
	}
}
