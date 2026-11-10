package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

const (
	// cacheImmutable is for /assets/*: Vite puts a content hash in their names, so a changed
	// file gets a new URL and the old one can be cached forever.
	cacheImmutable = "public, max-age=31536000, immutable"
	// cacheRevalidate is for everything else (index.html, favicon.svg), which keeps its URL
	// across deploys and must be checked each time.
	cacheRevalidate = "no-cache"
)

// spaHandler serves the SPA's files from dist, and index.html for any other path without a
// file extension, so client-side routes like /todos/1 load the app. Paths with an extension
// (a missing /assets/x.js) get a 404 instead of HTML. Directories are never listed.
func spaHandler(dist fs.FS) http.Handler {
	files := http.FileServerFS(dist)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if info, err := fs.Stat(dist, name); err == nil && !info.IsDir() {
			if strings.HasPrefix(name, "assets/") {
				w.Header().Set("Cache-Control", cacheImmutable)
			} else {
				w.Header().Set("Cache-Control", cacheRevalidate)
			}
			files.ServeHTTP(w, r)
			return
		}
		if path.Ext(name) != "" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", cacheRevalidate)
		http.ServeFileFS(w, r, dist, "index.html")
	})
}
