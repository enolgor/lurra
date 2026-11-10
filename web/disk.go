//go:build !embed

package web //nolint:revive // part of the app, not a library: no package doc needed

import (
	"io/fs"
	"log/slog"
	"os"
)

const distDir = "web/dist"

// Dist returns the built SPA from disk, relative to the working directory
func Dist() fs.FS {
	if _, err := os.Stat(distDir); err != nil {
		slog.Warn("SPA not found, UI will return 404", "dir", distDir)
	} else {
		slog.Info("serving SPA from disk", "dir", distDir)
	}
	return os.DirFS(distDir)
}
