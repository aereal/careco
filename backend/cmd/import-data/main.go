package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"careco/backend/log/attribute"
)

func main() { os.Exit(do()) }

func do() int {
	return exitCodeOf(run(context.Background()))
}

func run(ctx context.Context) error {
	entrypoint, err := build(ctx)
	if err != nil {
		return err
	}
	defer entrypoint.Shutdown(context.WithoutCancel(ctx))
	return entrypoint.Usecase.ImportData(ctx)
}

func exitCodeOf(err error) (exitCode int) {
	defer func() {
		level := slog.LevelDebug
		if exitCode > 0 {
			level = slog.LevelError
		}
		attrs := make([]slog.Attr, 0, 2)
		attrs = append(attrs, slog.Int("exit_code", exitCode))
		if err != nil {
			attrs = append(attrs, attribute.Error(err))
		}
		slog.LogAttrs(context.Background(), level, "application exited", attrs...)
	}()
	if err == nil {
		return 0
	}
	var hasExitCode interface{ ExitCode() int }
	if errors.As(err, &hasExitCode) {
		return hasExitCode.ExitCode()
	}
	return 1
}
