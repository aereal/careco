package o11y

import (
	"context"
	"log/slog"
	"slices"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type spanLogger struct{}

var _ sdktrace.SpanExporter = (*spanLogger)(nil)

func (*spanLogger) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	slog.DebugContext(ctx, "ExportSpans: start process spans", slog.Int("num", len(spans)))
	sortedSpans := slices.SortedFunc(slices.Values(spans), func(a, b sdktrace.ReadOnlySpan) int { return a.StartTime().Compare(b.StartTime()) })
	for _, span := range sortedSpans {
		if !span.SpanContext().IsValid() {
			slog.WarnContext(ctx, "invalid span",
				slog.String("name", span.Name()),
				slog.String("instrumentation.name", span.InstrumentationScope().Name),
			)
			continue
		}
		isRoot := !span.Parent().IsValid()
		switch {
		case !span.SpanContext().IsSampled():
			slog.DebugContext(ctx, "not-sampled span",
				slog.String("name", span.Name()),
				slog.Bool("is_root", isRoot),
				slog.String("trace_id", span.SpanContext().TraceID().String()),
				slog.String("span_id", span.SpanContext().SpanID().String()),
				slog.String("parent.trace_id", span.Parent().TraceID().String()),
				slog.String("parent.span_id", span.Parent().SpanID().String()),
				slog.String("instrumentation.name", span.InstrumentationScope().Name),
			)
		default:
			slog.DebugContext(ctx, "valid span",
				slog.String("name", span.Name()),
				slog.Bool("is_root", isRoot),
				slog.String("trace_id", span.SpanContext().TraceID().String()),
				slog.String("span_id", span.SpanContext().SpanID().String()),
				slog.String("parent.trace_id", span.Parent().TraceID().String()),
				slog.String("parent.span_id", span.Parent().SpanID().String()),
				slog.String("instrumentation.name", span.InstrumentationScope().Name),
			)
		}
	}
	return nil
}

func (*spanLogger) Shutdown(ctx context.Context) error {
	slog.DebugContext(ctx, "shutting down spanLogger")
	return nil
}
