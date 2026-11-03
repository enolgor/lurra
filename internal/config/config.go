package config

import (
	"errors"
	"io/fs"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int        `env:"PORT" envDefault:"8080"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"INFO"`
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, err
	}
	return env.ParseAs[Config]()
}
