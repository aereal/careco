package domain

import "context"

type TireCommand interface {
	AddTire(ctx context.Context, tire *Tire) (*Tire, error)
	RetireTire(ctx context.Context, id TireID) (*Tire, error)
}

type TireQuery interface {
	FindTires(ctx context.Context, includeRetired bool) ([]*Tire, error)
	FindTire(ctx context.Context, id TireID) (*Tire, error)
}
