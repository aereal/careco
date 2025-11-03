package internal

import (
	"careco/backend/web"
	"context"
	"log/slog"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ProvideEntrypoint(tp *sdktrace.TracerProvider, srv *web.Server) *Entrypoint {
	return &Entrypoint{
		Server: srv,
		tp:     tp,
	}
}

type Entrypoint struct {
	*web.Server

	tp *sdktrace.TracerProvider
}

func (e *Entrypoint) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := e.tp.Shutdown(ctx); err != nil {
		slog.WarnContext(ctx, "failed to shutdown OTel tracer provider")
	}
}
