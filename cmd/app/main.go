package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enolgor/lurra/internal/config"
	"github.com/enolgor/lurra/internal/controller"
	"github.com/enolgor/lurra/internal/logging"
	"github.com/enolgor/lurra/internal/server"
)

// version is set at build time via -ldflags
var version = "dev"

// shutdownTimeout is the time all services get to stop; keep it below Docker's 10s stop timeout.
const shutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("loading config", "error", err)
		os.Exit(1)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	slog.Info("starting", "version", version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // a second signal kills the process immediately

	c := controller.NewController()

	srv := server.New(cfg.Port, logger, c)
	slog.Info("listening", "addr", srv.Addr)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err := <-errCh:
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutting down server", "error", err)
	}
	// Close other services here, after the server stopped using them (e.g. db.Close(shutdownCtx)).

	slog.Info("stopped")
}
