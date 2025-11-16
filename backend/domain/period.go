package domain

func OpenEndpoint[T any](v T) Endpoint[T] { return Endpoint[T]{Value: v, Open: true} }

func ClosedEndpoint[T any](v T) Endpoint[T] { return Endpoint[T]{Value: v, Open: false} }

type Endpoint[T any] struct {
	Value T
	Open  bool
}

type Interval[T any] struct {
	Start, End Endpoint[T]
}

func EmptyInterval[T any]() Interval[T] {
	return Interval[T]{}
}
