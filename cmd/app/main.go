package main //nolint:revive // the app's entry point, not a library: no package doc needed

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enolgor/lurra/internal/api"
	"github.com/enolgor/lurra/internal/config"
	"github.com/enolgor/lurra/internal/logging"
	"github.com/enolgor/lurra/internal/server"
	"github.com/enolgor/lurra/internal/store/sqlite"
	"github.com/enolgor/lurra/web"
)

// version is set at build time via -ldflags
var version = "dev"

// shutdownTimeout is the time all services get to stop; keep it below Docker's 10s stop timeout.
const shutdownTimeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("exiting", "error", err)
		os.Exit(1)
	}
	slog.Info("stopped")
}

// run starts the app and blocks until a signal stops it (nil) or something fails (the error).
// It returns instead of calling os.Exit, so every deferred cleanup runs on every path.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger := slog.New(logging.WithRequestID(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))
	slog.SetDefault(logger)

	slog.Info("starting", "version", version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop) // a second signal kills the process immediately

	store, err := sqlite.Open(ctx, cfg.DBPath)
	if err != nil {
		return fmt.Errorf("opening store %s: %w", cfg.DBPath, err)
	}
	// Deferred calls run last-in first-out: this one runs after the server has shut down below,
	// once nothing uses the store any more.
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			slog.Error("closing store", "error", closeErr)
		}
	}()
	c := api.New(store)

	srv := server.New(cfg.Port, logger, c, web.Dist())
	slog.Info("listening", "addr", srv.Addr)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err := <-errCh:
		return fmt.Errorf("server stopped: %w", err)
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutting down server", "error", err)
	}
	return nil
}
