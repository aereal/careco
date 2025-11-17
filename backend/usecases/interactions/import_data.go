package interactions

import (
	"context"

	"careco/backend/usecases"

	"go.opentelemetry.io/otel/trace"
)

func ProvideImportData(tp trace.TracerProvider) *ImportData {
	return &ImportData{
		tracer: tp.Tracer("careco/backend/usecases/interactions.ImportData"),
	}
}

type ImportData struct {
	tracer trace.Tracer
}

var _ usecases.ImportData = (*ImportData)(nil)

func (u *ImportData) ImportData(ctx context.Context) (err error) {
	panic("TODO")
}
