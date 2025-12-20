package log

import (
	"context"
	"fmt"
	"log/slog"

	"careco/backend/o11y"

	"go.opentelemetry.io/otel/trace"
)

func injectGCPTraceAttrs(projectID string) middleware {
	return func(next handler) handler {
		return handlerFunc(func(ctx context.Context, record slog.Record) error {
			sc := trace.SpanContextFromContext(ctx)
			if !sc.IsValid() {
				return next.Handle(ctx, record)
			}
			// Google Cloud Logging形式のトレースフィールド
			// ref: https://cloud.google.com/trace/docs/trace-log-integration
			traceValue := fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID().String())
			record.AddAttrs(
				slog.String("logging.googleapis.com/trace", traceValue),
				slog.String("logging.googleapis.com/spanId", sc.SpanID().String()),
				slog.Bool("logging.googleapis.com/trace_sampled", sc.TraceFlags().IsSampled()),
			)
			return next.Handle(ctx, record)
		})
	}
}

func injectServiceVersion(svcVersion o11y.ServiceVersion) middleware {
	return func(next handler) handler {
		return handlerFunc(func(ctx context.Context, record slog.Record) error {
			record.AddAttrs(slog.GroupAttrs("svc", slog.String("version", string(svcVersion))))
			return next.Handle(ctx, record)
		})
	}
}

type handler interface {
	Handle(context.Context, slog.Record) error
}

var _ handler = (slog.Handler)(nil)

type handlerFunc func(context.Context, slog.Record) error

var _ handler = (handlerFunc)(nil)

func (f handlerFunc) Handle(ctx context.Context, record slog.Record) error { return f(ctx, record) }

type middleware func(handler) handler

func stack(edge slog.Handler, mws ...middleware) slog.Handler {
	return &stackedHandler{edge: edge, mws: mws}
}

type stackedHandler struct {
	edge slog.Handler
	mws  []middleware
}

var _ slog.Handler = (*stackedHandler)(nil)

func (h *stackedHandler) Handle(ctx context.Context, record slog.Record) error {
	var handler handler = h.edge
	for _, mw := range h.mws {
		handler = mw(handler)
	}
	return handler.Handle(ctx, record)
}

func (h *stackedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.edge.Enabled(ctx, level)
}

func (h *stackedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &stackedHandler{edge: h.edge.WithAttrs(attrs), mws: h.mws}
}

func (h *stackedHandler) WithGroup(name string) slog.Handler {
	return &stackedHandler{edge: h.edge.WithGroup(name), mws: h.mws}
}
