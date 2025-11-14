package tests

import (
	"fmt"
	"slices"

	"github.com/aereal/optional"
	"github.com/google/go-cmp/cmp"
)

func EquateOptional[T any]() cmp.Option {
	return cmp.Transformer(
		fmt.Sprintf("option%T", *new(T)),
		func(o optional.Option[T]) *T { return o.Ptr() },
	)
}

func Diff[T any](want, got T, opts ...cmp.Option) error {
	actOpts := slices.Clone(opts)
	if diff := cmp.Diff(want, got, actOpts...); diff != "" {
		return &MismatchError[T]{Diff: diff}
	}
	return nil
}

type MismatchError[T any] struct {
	Diff string
}

var _ error = (*MismatchError[any])(nil)

func (err *MismatchError[T]) Error() string {
	return "mismatch (-want, +got):\n" + err.Diff
}
