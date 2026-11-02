//go:build embed

package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed all:dist
var embedded embed.FS

// Dist returns the built SPA embedded in the binary
func Dist() fs.FS {
	dist, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	log.Print("serving SPA from embedded files")
	return dist
}
