package httpx_test

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enolgor/lurra/internal/httpx"
	"github.com/enolgor/lurra/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
)

// captureDefaultLog points the global slog default at a buffer for the rest of the test.
// The default is process-wide, so tests using this must not call t.Parallel.
func captureDefaultLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(logging.WithRequestID(slog.NewJSONHandler(&buf, nil))))
	t.Cleanup(func() { slog.SetDefault(old) })
	return &buf
}

// TestJSON_EncodeError checks that a value that can't be encoded is logged with the request ID.
// The status is already sent by then, so the client still gets it with an empty body.
func TestJSON_EncodeError(t *testing.T) {
	buf := captureDefaultLog(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, "abc"))
	rec := httptest.NewRecorder()

	httpx.JSON(rec, req, http.StatusOK, make(chan int)) // channels can't be encoded as JSON

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body)
	}

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("decode log line %q: %v", buf, err)
	}
	if got := m["level"]; got != "ERROR" {
		t.Errorf("level = %v, want %q", got, "ERROR")
	}
	if got := m["msg"]; got != "writing JSON response" {
		t.Errorf("msg = %v, want %q", got, "writing JSON response")
	}
	if _, ok := m["error"]; !ok {
		t.Error("error attribute missing")
	}
	if got := m["request_id"]; got != "abc" {
		t.Errorf("request_id = %v, want %q", got, "abc")
	}
}
