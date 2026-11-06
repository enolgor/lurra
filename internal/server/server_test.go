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

	"github.com/enolgor/lurra/internal/api"
	"github.com/enolgor/lurra/internal/model"
	"github.com/enolgor/lurra/internal/server"
	"github.com/enolgor/lurra/internal/store/sqlite"
)

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
	return server.New(0, logger, api.New(store)).Handler
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
