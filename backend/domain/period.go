package domain

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
