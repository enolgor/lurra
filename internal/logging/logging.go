package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

// New returns a text logger that adds the request ID to records logged with a
// request context, e.g. slog.InfoContext(ctx, ...).
func New(level slog.Level) *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(requestIDHandler{h})
}

type requestIDHandler struct{ slog.Handler }

func (h requestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := middleware.GetReqID(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h requestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return requestIDHandler{h.Handler.WithAttrs(attrs)}
}

func (h requestIDHandler) WithGroup(name string) slog.Handler {
	return requestIDHandler{h.Handler.WithGroup(name)}
}
