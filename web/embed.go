//go:build embed

package web //nolint:revive // shouldn't be imported by other modules

import (
	"embed"
	"io/fs"
	"log/slog"
)

//go:embed all:dist
var embedded embed.FS

// Dist returns the built SPA embedded in the binary
func Dist() fs.FS {
	dist, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	slog.Info("serving SPA from embedded files")
	return dist
}
