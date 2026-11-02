//go:build !embed

package web

import (
	"io/fs"
	"log"
	"os"
)

const distDir = "web/dist"

// Dist returns the built SPA from disk, relative to the working directory
func Dist() fs.FS {
	if _, err := os.Stat(distDir); err != nil {
		log.Printf("warning: SPA not found in %s, UI will return 404", distDir)
	} else {
		log.Printf("serving SPA from disk: %s", distDir)
	}
	return os.DirFS(distDir)
}
