package importdata

import (
	"context"
	"log/slog"
	"time"

	"careco/backend/log"
	"careco/backend/usecases"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ProvideEntrypoint(_ log.GlobalInstrumentationToken, tp *sdktrace.TracerProvider, importData usecases.ImportData) *Entrypoint {
	return &Entrypoint{
		Usecase: importData,
		tp:      tp,
	}
}

type Entrypoint struct {
	Usecase usecases.ImportData

	tp *sdktrace.TracerProvider
}

func (e *Entrypoint) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := e.tp.Shutdown(ctx); err != nil {
		slog.WarnContext(ctx, "failed to shutdown OTel tracer provider")
	}
}
