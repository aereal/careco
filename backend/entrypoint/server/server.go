package server

import (
	"context"
	"log/slog"
	"time"

	"careco/backend/infra/gcp"
	"careco/backend/log"
	"careco/backend/o11y"
	"careco/backend/web"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ProvideEntrypoint(_ log.GlobalInstrumentationToken, _ o11y.GlobalLoggerInstrumentationToken, tp *sdktrace.TracerProvider, srv *web.Server) *Entrypoint {
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

func BindGCPProject(v gcp.ProjectID) log.GoogleCloudProject { return log.GoogleCloudProject(v) }
