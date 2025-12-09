package log

import (
	"io"
	"log/slog"
	"os"

	"careco/backend/o11y"
)

type Output io.Writer

func ProvideStdoutOutput() Output { return os.Stdout }

func ProvideHandler(out Output, level slog.Level, svcVersion o11y.ServiceVersion) slog.Handler {
	return stack(
		slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}),
		injectOtelAttrs,
		injectServiceVersion(svcVersion),
		transformError,
	)
}

func ProvideJSONLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

type GlobalInstrumentationToken struct{}

func ProvideGlobalInstrumentation(logger *slog.Logger) GlobalInstrumentationToken {
	slog.SetDefault(logger)
	return GlobalInstrumentationToken{}
}
