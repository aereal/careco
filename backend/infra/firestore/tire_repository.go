package firestore

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"careco/backend/domain"
	"careco/backend/o11y/traceutils"

	"cloud.google.com/go/firestore"
	"github.com/aereal/optional"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const tiresCollection = "tires"

func ProvideTireRepository(tp trace.TracerProvider, client *firestore.Client) *TireRepository {
	return &TireRepository{
		tracer: tp.Tracer("careco/backend/infra/firestore.TireRepository"),
		client: client,
	}
}

type TireRepository struct {
	tracer trace.Tracer
	client *firestore.Client
}

var (
	_ domain.TireCommand = (*TireRepository)(nil)
	_ domain.TireQuery   = (*TireRepository)(nil)
)

func (r *TireRepository) AddTire(ctx context.Context, tire *domain.Tire) (_ *domain.Tire, err error) {
	ctx, span := r.tracer.Start(ctx, "AddTire")
	defer func() { traceutils.FinishSpan(span, err) }()

	docRef := r.client.Collection(tiresCollection).NewDoc()
	data := &dtoTire{
		Name:        tire.Name,
		PurchasedOn: tire.PurchasedOn.Ptr(),
		Retired:     tire.Retired,
	}
	if _, err := docRef.Create(ctx, data); err != nil {
		return nil, fmt.Errorf("firestore.DocumentRef.Create: %w", err)
	}
	return data.toDomain(docRef.ID), nil
}

func (r *TireRepository) RetireTire(ctx context.Context, id domain.TireID) (_ *domain.Tire, err error) {
	ctx, span := r.tracer.Start(ctx, "RetireTire")
	defer func() { traceutils.FinishSpan(span, err) }()

	docRef := r.client.Collection(tiresCollection).Doc(string(id))
	if _, err := docRef.Update(ctx, []firestore.Update{{Path: "Retired", Value: true}}); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrTireNotFound
		}
		return nil, fmt.Errorf("firestore.DocumentRef.Update: %w", err)
	}
	return r.FindTire(ctx, id)
}

func (r *TireRepository) FindTire(ctx context.Context, id domain.TireID) (_ *domain.Tire, err error) {
	ctx, span := r.tracer.Start(ctx, "FindTire")
	defer func() { traceutils.FinishSpan(span, err) }()

	doc, err := r.client.Collection(tiresCollection).Doc(string(id)).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, domain.ErrTireNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("firestore.DocumentRef.Get: %w", err)
	}
	dto := new(dtoTire)
	if err := doc.DataTo(dto); err != nil {
		return nil, fmt.Errorf("firestore.DocumentSnapshot.DataTo: %w", err)
	}
	return dto.toDomain(doc.Ref.ID), nil
}

func (r *TireRepository) FindTires(ctx context.Context, includeRetired bool) (_ []*domain.Tire, err error) {
	ctx, span := r.tracer.Start(ctx, "FindTires")
	defer func() { traceutils.FinishSpan(span, err) }()

	query := r.client.Collection(tiresCollection).Query
	if !includeRetired {
		query = query.WhereEntity(firestore.PropertyFilter{Path: "Retired", Operator: "==", Value: false})
	}
	docs := query.Documents(ctx)
	defer docs.Stop()
	tires := make([]*domain.Tire, 0)
	for {
		doc, err := docs.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("firestore.DocumentIterator.Next: %w", err)
		}
		dto := new(dtoTire)
		if err := doc.DataTo(dto); err != nil {
			return nil, fmt.Errorf("firestore.DocumentSnapshot.DataTo: %w", err)
		}
		tires = append(tires, dto.toDomain(doc.Ref.ID))
	}
	slices.SortStableFunc(tires, func(a, b *domain.Tire) int { return cmp.Compare(a.Name, b.Name) })
	return tires, nil
}

type dtoTire struct {
	Name        string
	PurchasedOn *time.Time
	Retired     bool
}

func (dto *dtoTire) toDomain(id string) *domain.Tire {
	return &domain.Tire{
		ID:          domain.TireID(id),
		Name:        dto.Name,
		PurchasedOn: optional.FromPtr(dto.PurchasedOn),
		Retired:     dto.Retired,
	}
}
