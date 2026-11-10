package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enolgor/lurra/internal/httpx"
)

func TestDecode(t *testing.T) {
	type payload struct {
		Title string `json:"title"`
	}
	tests := []struct {
		name       string
		body       string
		limit      int64 // max body size; 0 means no limit
		wantStatus int   // 0 means success
		want       payload
	}{
		{name: "valid", body: `{"title":"buy milk"}`, want: payload{Title: "buy milk"}},
		{name: "unknown field", body: `{"titel":"buy milk"}`, wantStatus: http.StatusBadRequest},
		{name: "malformed", body: `{"title":`, wantStatus: http.StatusBadRequest},
		{name: "trailing data", body: `{"title":"a"} {"title":"b"}`, wantStatus: http.StatusBadRequest},
		{name: "too large", body: `{"title":"buy milk"}`, limit: 5, wantStatus: http.StatusRequestEntityTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.limit > 0 {
				req.Body = http.MaxBytesReader(httptest.NewRecorder(), req.Body, tt.limit)
			}

			var got payload
			err := httpx.Decode(req, &got)

			if tt.wantStatus == 0 {
				if err != nil {
					t.Fatalf("Decode() error = %v, want nil", err)
				}
				if got != tt.want {
					t.Errorf("decoded = %+v, want %+v", got, tt.want)
				}
				return
			}
			se, ok := errors.AsType[*httpx.StatusError](err)
			if !ok {
				t.Fatalf("Decode() error = %v (%T), want a *StatusError", err, err)
			}
			if se.Status != tt.wantStatus {
				t.Errorf("status = %d, want %d (error: %v)", se.Status, tt.wantStatus, err)
			}
		})
	}
}
