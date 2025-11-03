package main

import (
	"context"
	"errors"
	"os"
)

func main() { os.Exit(do()) }

func do() int {
	return exitCodeOf(run(context.Background()))
}

func run(ctx context.Context) error {
	srv := build()
	return srv.Start(ctx)
}

func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	var hasExitCode interface{ ExitCode() int }
	if errors.As(err, &hasExitCode) {
		return hasExitCode.ExitCode()
	}
	return 1
}
