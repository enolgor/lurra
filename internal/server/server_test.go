package server_test

import (
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/enolgor/lurra/internal/api"
	"github.com/enolgor/lurra/internal/model"
	"github.com/enolgor/lurra/internal/server"
	"github.com/enolgor/lurra/internal/store/sqlite"
)

// testDist stands in for web/dist.
var testDist = fstest.MapFS{
	"index.html":      {Data: []byte("<!doctype html>app")},
	"favicon.svg":     {Data: []byte("<svg></svg>")},
	"assets/app-1.js": {Data: []byte("console.log(1)")},
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	store, err := sqlite.Open(t.Context(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})
	logger := slog.New(slog.DiscardHandler)
	return server.New(0, logger, api.New(store), testDist).Handler
}

func doRequest(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// TestAPI_Errors covers only what needs the full server: the {id} route parameter, the
// real store's not-found error and the RequestSize middleware. The api tests cover the rest.
func TestAPI_Errors(t *testing.T) {
	// just over maxRequestBody (1 MB) in server.go
	bigBody := `{"title":"` + strings.Repeat("a", 1<<20) + `"}`
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{name: "invalid id", method: http.MethodGet, path: "/api/v1/todos/abc", wantStatus: http.StatusBadRequest},
		{name: "not found", method: http.MethodGet, path: "/api/v1/todos/999", wantStatus: http.StatusNotFound},
		{name: "body too large", method: http.MethodPost, path: "/api/v1/todos", body: bigBody, wantStatus: http.StatusRequestEntityTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t)
			rec := doRequest(t, h, tt.method, tt.path, tt.body)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}

func TestTodoLifecycle(t *testing.T) {
	h := newTestHandler(t)

	createRec := doRequest(t, h, http.MethodPost, "/api/v1/todos", `{"title": "buy milk"}`)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want %d; body: %s", createRec.Code, http.StatusCreated, createRec.Body)
	}
	var created model.Todo
	if err := json.UnmarshalRead(createRec.Body, &created); err != nil {
		t.Fatalf("create: decode body: %v", err)
	}
	if created.ID == 0 {
		t.Errorf("create: ID = 0, want non-zero")
	}
	if created.Title != "buy milk" {
		t.Errorf("create: Title = %q, want %q", created.Title, "buy milk")
	}
	if created.Done {
		t.Errorf("create: Done = true, want false")
	}

	listRec := doRequest(t, h, http.MethodGet, "/api/v1/todos", "")
	if listRec.Code != http.StatusOK {
		t.Fatalf("list: status = %d, want %d; body: %s", listRec.Code, http.StatusOK, listRec.Body)
	}
	var list []model.Todo
	if err := json.UnmarshalRead(listRec.Body, &list); err != nil {
		t.Fatalf("list: decode body: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("list: len = %d, want 1; list: %+v", len(list), list)
	}
	if list[0] != created {
		t.Errorf("list: list[0] = %+v, want %+v", list[0], created)
	}

	path := fmt.Sprintf("/api/v1/todos/%d", created.ID)

	readRec := doRequest(t, h, http.MethodGet, path, "")
	if readRec.Code != http.StatusOK {
		t.Fatalf("read: status = %d, want %d; body: %s", readRec.Code, http.StatusOK, readRec.Body)
	}
	var read model.Todo
	if err := json.UnmarshalRead(readRec.Body, &read); err != nil {
		t.Fatalf("read: decode body: %v", err)
	}
	if read != created {
		t.Errorf("read: read = %+v, want %+v", read, created)
	}

	want := created
	want.Done = true
	updateRec := doRequest(t, h, http.MethodPatch, path, `{"done":true}`)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("update: status = %d, want %d; body: %s", updateRec.Code, http.StatusOK, updateRec.Body)
	}
	var updated model.Todo
	if err := json.UnmarshalRead(updateRec.Body, &updated); err != nil {
		t.Fatalf("update: decode body: %v", err)
	}
	if updated != want {
		t.Errorf("update: updated = %+v, want %+v", updated, want)
	}

	deleteRec := doRequest(t, h, http.MethodDelete, path, "")
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete: status = %d, want %d; body: %s", deleteRec.Code, http.StatusNoContent, deleteRec.Body)
	}

	goneRec := doRequest(t, h, http.MethodGet, path, "")
	if goneRec.Code != http.StatusNotFound {
		t.Fatalf("read after delete: status = %d, want %d; body: %s", goneRec.Code, http.StatusNotFound, goneRec.Body)
	}
}

// TestCrossOriginProtection checks that other sites can't change state through the user's
// browser (CSRF), while the SPA (same origin), non-browser clients and reads still work.
func TestCrossOriginProtection(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		headers    map[string]string
		wantStatus int
	}{
		{name: "cross-site POST", method: http.MethodPost, headers: map[string]string{"Sec-Fetch-Site": "cross-site"}, wantStatus: http.StatusForbidden},
		{name: "same-site POST", method: http.MethodPost, headers: map[string]string{"Sec-Fetch-Site": "same-site"}, wantStatus: http.StatusForbidden},
		{name: "other Origin POST", method: http.MethodPost, headers: map[string]string{"Origin": "https://evil.example"}, wantStatus: http.StatusForbidden},
		{name: "same-origin POST", method: http.MethodPost, headers: map[string]string{"Sec-Fetch-Site": "same-origin"}, wantStatus: http.StatusCreated},
		{name: "matching Origin POST", method: http.MethodPost, headers: map[string]string{"Origin": "http://example.com"}, wantStatus: http.StatusCreated},
		{name: "POST without browser headers", method: http.MethodPost, wantStatus: http.StatusCreated},
		{name: "cross-site GET", method: http.MethodGet, headers: map[string]string{"Sec-Fetch-Site": "cross-site"}, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t)
			// httptest requests have Host example.com
			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/api/v1/todos", strings.NewReader(`{"title":"x"}`))
			for name, value := range tt.headers {
				req.Header.Set(name, value)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}

func TestSPA(t *testing.T) {
	const (
		index      = "<!doctype html>app"
		immutable  = "public, max-age=31536000, immutable"
		revalidate = "no-cache"
	)
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string // checked when set
		wantCache  string // Cache-Control
	}{
		{name: "root serves index.html", path: "/", wantStatus: http.StatusOK, wantBody: index, wantCache: revalidate},
		{name: "client route serves index.html", path: "/todos/1", wantStatus: http.StatusOK, wantBody: index, wantCache: revalidate},
		{name: "directory serves index.html, not a listing", path: "/assets/", wantStatus: http.StatusOK, wantBody: index, wantCache: revalidate},
		{name: "hashed asset", path: "/assets/app-1.js", wantStatus: http.StatusOK, wantBody: "console.log(1)", wantCache: immutable},
		{name: "unhashed top-level file", path: "/favicon.svg", wantStatus: http.StatusOK, wantBody: "<svg></svg>", wantCache: revalidate},
		{name: "missing asset is not cached", path: "/assets/nope.js", wantStatus: http.StatusNotFound},
		{name: "unknown API route", path: "/api/v1/nope", wantStatus: http.StatusNotFound},
		{name: "unknown API version", path: "/api/v2/todos", wantStatus: http.StatusNotFound},
	}

	h := newTestHandler(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, h, http.MethodGet, tt.path, "")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body, tt.wantBody)
			}
			if got := rec.Header().Get("Cache-Control"); got != tt.wantCache {
				t.Errorf("Cache-Control = %q, want %q", got, tt.wantCache)
			}
			if tt.wantStatus == http.StatusNotFound && strings.Contains(rec.Body.String(), "<!doctype") {
				t.Errorf("404 returned the SPA: %s", rec.Body)
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	want := map[string]string{
		"X-Content-Type-Options":     "nosniff",
		"X-Frame-Options":            "DENY",
		"Referrer-Policy":            "strict-origin-when-cross-origin",
		"Cross-Origin-Opener-Policy": "same-origin",
	}
	// CSP directives that must be present; the full policy is in headers.go.
	wantCSP := []string{"default-src 'self'", "script-src 'self'", "object-src 'none'", "frame-ancestors 'none'"}

	h := newTestHandler(t)
	for _, path := range []string{"/", "/todos/1", "/api/v1/todos", "/api/v1/todos/999"} {
		t.Run(path, func(t *testing.T) {
			rec := doRequest(t, h, http.MethodGet, path, "")
			for name, value := range want {
				if got := rec.Header().Get(name); got != value {
					t.Errorf("%s = %q, want %q", name, got, value)
				}
			}
			csp := rec.Header().Get("Content-Security-Policy")
			for _, directive := range wantCSP {
				if !strings.Contains(csp, directive) {
					t.Errorf("Content-Security-Policy = %q, missing %q", csp, directive)
				}
			}
			if strings.Contains(csp, "script-src 'self' 'unsafe-inline'") {
				t.Errorf("Content-Security-Policy allows inline scripts: %q", csp)
			}
		})
	}
}
