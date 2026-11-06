package logging_test

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/enolgor/lurra/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
)

func newTestLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	return slog.New(logging.WithRequestID(slog.NewJSONHandler(&buf, nil))), &buf
}

func decode(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("decode log line %q: %v", buf, err)
	}
	return m
}

func TestWithRequestID(t *testing.T) {
	t.Run("no context", func(t *testing.T) {
		logger, buf := newTestLogger(t)
		logger.Info("some message")
		m := decode(t, buf)
		if v, ok := m["request_id"]; ok {
			t.Errorf("request_id = %v, want absent", v)
		}
	})

	t.Run("context with id", func(t *testing.T) {
		logger, buf := newTestLogger(t)
		want := "abc"
		ctx := context.WithValue(t.Context(), middleware.RequestIDKey, want)
		logger.InfoContext(ctx, "some message")
		m := decode(t, buf)
		if got := m["request_id"]; got != want {
			t.Errorf("request_id = %v, want %q", got, want)
		}
	})
}

func TestWithRequestID_Slogtest(t *testing.T) {
	var buf bytes.Buffer
	slogtest.Run(t,
		func(t *testing.T) slog.Handler {
			buf.Reset()
			return logging.WithRequestID(slog.NewJSONHandler(&buf, nil))
		},
		func(t *testing.T) map[string]any {
			return decode(t, &buf)
		},
	)
}

// TestWithRequestID_WithGroup pins the documented limitation: request_id lands inside an
// open group, not at the top level. If you fix that, update this test and the WithRequestID doc.
func TestWithRequestID_WithGroup(t *testing.T) {
	logger, buf := newTestLogger(t)
	ctx := context.WithValue(t.Context(), middleware.RequestIDKey, "abc")
	logger.WithGroup("req").InfoContext(ctx, "hello", "k", "v")
	m := decode(t, buf)

	if got, ok := m["request_id"]; ok {
		t.Errorf("request_id = %v, want absent at the top level (it goes inside the group)", got)
	}

	group, ok := m["req"].(map[string]any)
	if !ok {
		t.Fatalf("req = %#v, want a group", m["req"])
	}
	if got := group["k"]; got != "v" {
		t.Errorf("req.k = %v, want %q", got, "v")
	}
	if got := group["request_id"]; got != "abc" {
		t.Errorf("req.request_id = %v, want %q", got, "abc")
	}
}
