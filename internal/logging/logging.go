package logging

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

// WithRequestID wraps h so it adds the request ID to records logged with a request
// context, e.g. slog.InfoContext(ctx, ...). The ID is a record attribute, so it ends up
// inside any group opened with WithGroup (e.g. req.request_id) rather than at the top level.
func WithRequestID(h slog.Handler) slog.Handler {
	return requestIDHandler{h}
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
