package httpx_test

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enolgor/lurra/internal/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

func TestWithStatus(t *testing.T) {
	base := errors.New("title is empty")
	err := httpx.WithStatus(http.StatusBadRequest, base)

	if got := err.Error(); got != base.Error() {
		t.Errorf("Error() = %q, want %q", got, base.Error())
	}
	if !errors.Is(err, base) {
		t.Error("errors.Is(err, base) = false, want true (Unwrap must return the wrapped error)")
	}
	se, ok := errors.AsType[*httpx.StatusError](err)
	if !ok {
		t.Fatalf("errors.AsType[*StatusError] failed for %T", err)
	}
	if se.Status != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", se.Status, http.StatusBadRequest)
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "status error",
			err:        httpx.WithStatus(http.StatusNotFound, errors.New("todo id=9 is not found")),
			wantStatus: http.StatusNotFound,
			wantMsg:    "todo id=9 is not found",
		},
		{
			// The status is found through wrapping; the message is the StatusError's own,
			// without the wrapping context.
			name:       "wrapped status error",
			err:        fmt.Errorf("handler: %w", httpx.WithStatus(http.StatusBadRequest, errors.New("title is empty"))),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "title is empty",
		},
		{
			// Errors without a status must not leak their message to the client.
			name:       "plain error",
			err:        errors.New("disk is on fire"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    http.StatusText(http.StatusInternalServerError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, "abc"))
			rec := httptest.NewRecorder()

			httpx.NewError(rec, req, tt.err)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json")
			}
			var body httpx.HTTPError
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body %q: %v", rec.Body, err)
			}
			want := httpx.HTTPError{Code: tt.wantStatus, Message: tt.wantMsg, RequestID: "abc"}
			if body != want {
				t.Errorf("body = %+v, want %+v", body, want)
			}
		})
	}
}

// TestNewError_RequestLog checks that the full error, even one hidden from the client,
// ends up on the request's log line written by httplog.
func TestNewError_RequestLog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	h := httplog.RequestLogger(logger, &httplog.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.NewError(w, r, errors.New("disk is on fire"))
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))

	if strings.Contains(rec.Body.String(), "disk is on fire") {
		t.Errorf("body = %q, must not contain the internal error", rec.Body)
	}
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("decode log line %q: %v", buf.String(), err)
	}
	if got := m["error"]; got != "disk is on fire" {
		t.Errorf("logged error = %v, want %q", got, "disk is on fire")
	}
}
