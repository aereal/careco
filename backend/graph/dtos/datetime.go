package dtos

import (
	"context"
	"io"
	"strconv"
	"time"

	"careco/backend/types"

	"github.com/99designs/gqlgen/graphql"
)

var layout = "2006-01-02T15:04:05.999-07:00"

func MarshalDateTime(t time.Time) graphql.ContextMarshaler {
	return graphql.ContextWriterFunc(func(_ context.Context, w io.Writer) error {
		s := strconv.Quote(t.Format(layout))
		_, err := w.Write([]byte(s))
		return err
	})
}

func UnmarshalDateTime(_ context.Context, v any) (time.Time, error) {
	s, err := types.Cast[string](v)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, s)
}
