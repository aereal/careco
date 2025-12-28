package log

import (
	"io"
	"log/slog"
	"os"

	"careco/backend/o11y"
)

type Output io.Writer

func ProvideStdoutOutput() Output { return os.Stdout }

func ProvideJSONLogger(out Output, level slog.Level, svcVersion o11y.ServiceVersion, gcpProject GoogleCloudProject) *slog.Logger {
	handler := stack(
		slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}),
		injectOtelAttrs(gcpProject),
		injectServiceVersion(svcVersion),
	)
	return slog.New(handler)
}

type GlobalInstrumentationToken struct{}

func ProvideGlobalInstrumentation(logger *slog.Logger) GlobalInstrumentationToken {
	slog.SetDefault(logger)
	return GlobalInstrumentationToken{}
}
