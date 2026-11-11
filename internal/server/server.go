package server

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/enolgor/lurra/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

const (
	maxRequestBody    = 1 << 20 // 1 MB
	compressLevel     = 5
	handlerTimeout    = 30 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 45 * time.Second // must exceed handlerTimeout
	idleTimeout       = 120 * time.Second
)

// newRouter returns the HTTP handler with the middleware stack, the API and the SPA in dist.
func newRouter(logger *slog.Logger, c *api.Controller, dist fs.FS) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping")) // before the logger, so health checks are not logged
	r.Use(middleware.RequestID)
	r.Use(securityHeaders)
	// Logs every request and recovers panics (replaces chi's Logger and Recoverer).
	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level:         slog.LevelDebug, // filtering is left to the logger's level
		Schema:        httplog.SchemaECS.Concise(true),
		RecoverPanics: true,
	}))
	// CSRF: rejects cross-origin POST/PATCH/DELETE with 403 (Sec-Fetch-Site or Origin header).
	// After the logger, so rejections are logged.
	r.Use(http.NewCrossOriginProtection().Handler)
	r.Use(middleware.Compress(compressLevel))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.RequestSize(maxRequestBody))
		r.Use(middleware.Timeout(handlerTimeout))
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", c.ListTodos)
			r.Post("/", c.CreateTodo)
			r.Get("/{id}", c.ReadTodo)
			r.Patch("/{id}", c.UpdateTodo)
			r.Delete("/{id}", c.DeleteTodo)
		})
	})
	r.Handle("/api/*", http.NotFoundHandler()) // other API versions: 404, not the SPA
	r.Handle("/*", spaHandler(dist))

	return r
}

// New returns the HTTP server on port, with the router and timeouts set. dist holds the
// built SPA.
func New(port int, logger *slog.Logger, c *api.Controller, dist fs.FS) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           newRouter(logger, c, dist),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
