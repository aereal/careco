package log

import (
	"careco/backend/log/attribute"
	"careco/backend/o11y"
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func injectOtelAttrs(next handler) handler {
	return handlerFunc(func(ctx context.Context, record slog.Record) error {
		sc := trace.SpanContextFromContext(ctx)
		if !sc.IsValid() {
			return next.Handle(ctx, record)
		}
		attrTraceID := slog.String("trace_id", sc.TraceID().String())
		attrSpanID := slog.String("span_id", sc.SpanID().String())
		newRecord := mergeRecordAttrs(record, mergeGroupAttrChildren("otel", attrTraceID, attrSpanID))
		return next.Handle(ctx, newRecord)
	})
}

func injectServiceVersion(svcVersion o11y.ServiceVersion) middleware {
	return func(next handler) handler {
		return handlerFunc(func(ctx context.Context, record slog.Record) error {
			record.AddAttrs(slog.GroupAttrs("svc", slog.String("version", string(svcVersion))))
			return next.Handle(ctx, record)
		})
	}
}

func transformError(next handler) handler {
	return handlerFunc(func(ctx context.Context, record slog.Record) error {
		newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
		for a := range record.Attrs {
			if a.Key == attribute.KeyError {
				if ev, ok := a.Value.Any().(*attribute.ErrorValue); ok {
					ga := slog.GroupAttrs(
						attribute.KeyError,
						slog.String("type", fmt.Sprintf("%T", ev.Err)),
						slog.String("msg", ev.Err.Error()),
					)
					newRecord.AddAttrs(ga)
				}
			} else {
				newRecord.AddAttrs(a)
			}
		}
		return next.Handle(ctx, newRecord)
	})
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
