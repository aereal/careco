package dtos

import (
	"context"
	"io"
	"iter"
	"maps"
	"strconv"
	"strings"
	"time"

	"careco/backend/types"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalMonth(m time.Month) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		if m < time.January || m > time.December {
			return ErrOutOfMonthRange
		}
		_, err := w.Write([]byte(strconv.Quote(strings.ToUpper(m.String()))))
		return err
	})
}

func UnmarshalMonth(_ context.Context, v any) (time.Month, error) {
	s, err := types.Cast[string](v)
	if err != nil {
		return 0, err
	}
	month, ok := monthByName[s]
	if !ok {
		return 0, &UnknownMonthError{Value: s}
	}
	return month, nil
}

var monthByName = maps.Collect(monthAndName())

func monthAndName() iter.Seq2[string, time.Month] {
	return func(yield func(string, time.Month) bool) {
		for i := range 12 {
			m := time.Month(i + 1)
			if !yield(m.String(), m) {
				return
			}
		}
	}
}
