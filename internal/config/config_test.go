package config_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/enolgor/lurra/internal/config"
)

var envKeys = []string{
	"PORT",
	"LOG_LEVEL",
	"DB_PATH",
}

func isolate(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	for _, key := range envKeys {
		t.Setenv(key, "")                        // registers a cleanup that restores the original value
		if err := os.Unsetenv(key); err != nil { // actually removes it, so godotenv sees it as missing
			t.Fatalf("unset %s: %v", key, err)
		}
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		dotenv  string
		want    config.Config
		wantErr bool
	}{
		{
			name: "defaults",
			want: config.Config{Port: 8080, LogLevel: slog.LevelInfo, DBPath: "data/app.db"},
		},
		{
			name: "from env",
			env:  map[string]string{"PORT": "7000", "LOG_LEVEL": "debug", "DB_PATH": "/tmp/x.db"},
			want: config.Config{Port: 7000, LogLevel: slog.LevelDebug, DBPath: "/tmp/x.db"},
		},
		{
			name:   "from .env",
			dotenv: `PORT=9000`,
			want:   config.Config{Port: 9000, LogLevel: slog.LevelInfo, DBPath: "data/app.db"},
		},
		{
			name:   "env overrides .env",
			env:    map[string]string{"PORT": "7000"},
			dotenv: `PORT=9000`,
			want:   config.Config{Port: 7000, LogLevel: slog.LevelInfo, DBPath: "data/app.db"},
		},
		{
			name:    "invalid port",
			env:     map[string]string{"PORT": "abc"},
			wantErr: true,
		},
		{
			name:    "invalid level",
			env:     map[string]string{"LOG_LEVEL": "loud"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			if tt.dotenv != "" {
				if err := os.WriteFile(".env", []byte(tt.dotenv), 0o600); err != nil {
					t.Fatalf("write .env: %v", err)
				}
			}
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("config = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLoad_DotenvUnreadable(t *testing.T) {
	isolate(t)
	if err := os.Mkdir(".env", 0o755); err != nil {
		t.Fatalf("mkdir .env: %v", err)
	}
	_, err := config.Load()
	if err == nil {
		t.Errorf("Load() error = nil, want an error")
	}
}
