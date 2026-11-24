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
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTireRepository_AddTire(t *testing.T) {
	t.Parallel()

	purchasedOn := time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name  string
		input *domain.Tire
		want  *domain.Tire
	}{
		{
			name:  "ok/with purchase date",
			input: &domain.Tire{Name: "summer", PurchasedOn: optional.Some(purchasedOn)},
			want:  &domain.Tire{Name: "summer", PurchasedOn: optional.Some(purchasedOn)},
		},
		{
			name:  "ok/without purchase date",
			input: &domain.Tire{Name: "studless"},
			want:  &domain.Tire{Name: "studless"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := test.BuildTireRepository(t)
			if err != nil {
				t.Fatal(err)
			}
			ctx := t.Context()
			added, err := r.AddTire(ctx, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if added.ID == "" {
				t.Fatal("added tire has no ID")
			}
			if err := tests.Diff(tc.want, added, ignoreTireID(), tests.EquateOptional[time.Time]()); err != nil {
				t.Errorf("AddTire: %s", err)
			}
			found, err := r.FindTire(ctx, added.ID)
			if err != nil {
				t.Fatal(err)
			}
			if err := tests.Diff(added, found, tests.EquateOptional[time.Time]()); err != nil {
				t.Errorf("FindTire: %s", err)
			}
		})
	}
}

func TestTireRepository_FindTire(t *testing.T) {
	t.Parallel()

	r, err := test.BuildTireRepository(t)
	if err != nil {
		t.Fatal(err)
	}
	_, gotErr := r.FindTire(t.Context(), domain.TireID("not-exist"))
	tests.AssertsErrors(t, domain.ErrTireNotFound, gotErr)
}

func TestTireRepository_FindTires(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		includeRetired bool
		want           []*domain.Tire
	}{
		{
			name:           "exclude retired",
			includeRetired: false,
			want: []*domain.Tire{
				{Name: "a"},
				{Name: "c"},
			},
		},
		{
			name:           "include retired",
			includeRetired: true,
			want: []*domain.Tire{
				{Name: "a"},
				{Name: "b", Retired: true},
				{Name: "c"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := test.BuildTireRepository(t)
			if err != nil {
				t.Fatal(err)
			}
			ctx := t.Context()
			if err := prepareTires(ctx, r); err != nil {
				t.Fatal(err)
			}
			got, err := r.FindTires(ctx, tc.includeRetired)
			if err != nil {
				t.Fatal(err)
			}
			if err := tests.Diff(tc.want, got, ignoreTireID(), tests.EquateOptional[time.Time]()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTireRepository_RetireTire(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		r, err := test.BuildTireRepository(t)
		if err != nil {
			t.Fatal(err)
		}
		ctx := t.Context()
		added, err := r.AddTire(ctx, &domain.Tire{Name: "summer"})
		if err != nil {
			t.Fatal(err)
		}
		retired, err := r.RetireTire(ctx, added.ID)
		if err != nil {
			t.Fatal(err)
		}
		want := &domain.Tire{ID: added.ID, Name: "summer", Retired: true}
		if err := tests.Diff(want, retired, tests.EquateOptional[time.Time]()); err != nil {
			t.Errorf("RetireTire: %s", err)
		}
		found, err := r.FindTire(ctx, added.ID)
		if err != nil {
			t.Fatal(err)
		}
		if err := tests.Diff(want, found, tests.EquateOptional[time.Time]()); err != nil {
			t.Errorf("FindTire: %s", err)
		}
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		r, err := test.BuildTireRepository(t)
		if err != nil {
			t.Fatal(err)
		}
		_, gotErr := r.RetireTire(t.Context(), domain.TireID("not-exist"))
		tests.AssertsErrors(t, domain.ErrTireNotFound, gotErr)
	})
}

func prepareTires(ctx context.Context, r *firestore.TireRepository) error {
	var errs error
	for _, name := range []string{"c", "a"} {
		if _, err := r.AddTire(ctx, &domain.Tire{Name: name}); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	b, err := r.AddTire(ctx, &domain.Tire{Name: "b"})
	if err != nil {
		return errors.Join(errs, err)
	}
	if _, err := r.RetireTire(ctx, b.ID); err != nil {
		errs = errors.Join(errs, err)
	}
	return errs
}

func ignoreTireID() cmp.Option {
	return cmpopts.IgnoreFields(domain.Tire{}, "ID")
}
