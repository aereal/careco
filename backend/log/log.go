package log

import (
	"io"
	"log/slog"
	"os"

	"careco/backend/o11y"
)

type Output io.Writer

func ProvideStdoutOutput() Output { return os.Stdout }

func ProvideHandler(out Output, level slog.Level, svcVersion o11y.ServiceVersion, gcpProject GoogleCloudProject) slog.Handler {
	return stack(
		slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}),
		injectOtelAttrs(gcpProject),
		injectServiceVersion(svcVersion),
	)
}

type GlobalInstrumentationToken struct{}

func ProvideGlobalInstrumentation(h slog.Handler) GlobalInstrumentationToken {
	slog.SetDefault(slog.New(h))
	return GlobalInstrumentationToken{}
}
