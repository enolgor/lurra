package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/enolgor/lurra/internal/config"
	"github.com/enolgor/lurra/internal/controller"
	"github.com/enolgor/lurra/internal/logging"
	"github.com/enolgor/lurra/internal/server"
)

// version is set at build time via -ldflags
var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("loading config", "error", err)
		os.Exit(1)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	slog.Info("starting", "version", version)

	c := controller.NewController()

	if err := server.Run(fmt.Sprintf(":%d", cfg.Port), server.NewRouter(logger, c)); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
