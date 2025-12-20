package log

import (
	"io"
	"log/slog"
	"os"

	"careco/backend/infra/gcp"
	"careco/backend/o11y"
)

type Output io.Writer

func ProvideStdoutOutput() Output { return os.Stdout }

func ProvideLogger(handler slog.Handler) *slog.Logger { return slog.New(handler) }

func ProvideJSONHandler(out Output, level slog.Level, svcVersion o11y.ServiceVersion, projectID gcp.ProjectID) slog.Handler {
	return stack(
		slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}),
		injectGCPTraceAttrs(string(projectID)),
		injectServiceVersion(svcVersion),
	)
}

type GlobalInstrumentationToken struct{}

func ProvideGlobalInstrumentation(logger *slog.Logger) GlobalInstrumentationToken {
	slog.SetDefault(logger)
	return GlobalInstrumentationToken{}
}
