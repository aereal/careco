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
)

func TestDrivingRecordRepository_RecordDrivingRecord(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		record  *domain.DrivingRecord
		wantErr error
	}{
		{
			name: "ok",
			record: &domain.DrivingRecord{
				Date:               time.Now(),
				DistanceKilometers: 123,
			},
			wantErr: nil,
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
			gotErr := r.RecordDrivingRecord(ctx, tc.record)
			tests.AssertsErrors(t, tc.wantErr, gotErr)
		})
	}
}

func TestDrivingRecordRepository_FindRecentRecords(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		first   int
		want    []*domain.DrivingRecord
		wantErr error
		prepare func(ctx context.Context, r *firestore.DrivingRecordRepository) error
	}{
		{
			name:  "ok",
			first: 3,
			want: []*domain.DrivingRecord{
				{
					DistanceKilometers: 2,
					Date:               time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
				},
				{
					DistanceKilometers: 3,
					Date:               time.Date(2025, time.January, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					DistanceKilometers: 1,
					Date:               time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
			prepare: func(ctx context.Context, r *firestore.DrivingRecordRepository) error {
				return errors.Join(
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 1, Date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 2, Date: time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC)}),
					r.RecordDrivingRecord(ctx, &domain.DrivingRecord{DistanceKilometers: 3, Date: time.Date(2025, time.January, 4, 0, 0, 0, 0, time.UTC)}),
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
			got, gotErr := r.FindRecentRecords(ctx, tc.first)
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
				Start: domain.Endpoint[time.Time]{Value: baseDate, Open: true},
				End:   domain.Endpoint[time.Time]{Value: nextMonth, Open: false},
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
				Start: domain.Endpoint[time.Time]{Value: baseDate, Open: true},
				End:   domain.Endpoint[time.Time]{Value: nextMonth, Open: false},
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
			got, gotErr := r.FindRecordsInPeriod(ctx, tc.interval, tc.direction)
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
