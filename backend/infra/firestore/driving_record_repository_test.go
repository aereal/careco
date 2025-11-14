package firestore_test

import (
	"testing"
	"time"

	"careco/backend/domain"
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
