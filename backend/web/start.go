package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func start(ctx context.Context, srv *http.Server, shutdownGrace time.Duration) error {
	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-signalCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(signalCtx), shutdownGrace)
		defer cancel()
		slog.InfoContext(shutdownCtx, "shutting down server", slog.Duration("grace", shutdownGrace))
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.WarnContext(shutdownCtx, "server has stopped inelegantly")
		}
	}()

	slog.InfoContext(ctx, "starting a server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
