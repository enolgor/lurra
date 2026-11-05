package api_test

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enolgor/lurra/internal/api"
	"github.com/enolgor/lurra/internal/httpx"
	"github.com/enolgor/lurra/internal/model"
	"github.com/go-chi/chi/v5"
)

func newRequest(t *testing.T, id, body string) *http.Request {
	t.Helper()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", strings.NewReader(body))
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

type failingStore struct {
	model.TodoStore
	err error
}

func (s failingStore) ListTodos(_ context.Context) ([]model.Todo, error) {
	return nil, s.err
}

func (s failingStore) CreateTodo(_ context.Context, _ string) (model.Todo, error) {
	return model.Todo{}, s.err
}

func (s failingStore) ReadTodo(_ context.Context, _ int) (model.Todo, error) {
	return model.Todo{}, s.err
}

func (s failingStore) UpdateTodo(_ context.Context, _ int, _ model.UpdateTodo) (model.Todo, error) {
	return model.Todo{}, s.err
}

func (s failingStore) DeleteTodo(_ context.Context, _ int) error {
	return s.err
}

func TestController_StoreError(t *testing.T) {
	secret := "disk is on fire"
	c := api.New(failingStore{err: errors.New(secret)})

	tests := []struct {
		name    string
		handler http.HandlerFunc
		id      string
		body    string
	}{
		{name: "ListTodos", handler: c.ListTodos},
		{name: "CreateTodo", handler: c.CreateTodo, body: `{"title":"buy milk"}`},
		{name: "ReadTodo", handler: c.ReadTodo, id: "1"},
		{name: "UpdateTodo", handler: c.UpdateTodo, id: "1", body: `{"done":true}`},
		{name: "DeleteTodo", handler: c.DeleteTodo, id: "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(t, tt.id, tt.body)
			rec := httptest.NewRecorder()
			tt.handler(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusInternalServerError, rec.Body)
			}
			if strings.Contains(rec.Body.String(), secret) {
				t.Errorf("body = %s, want absence of %q", rec.Body, secret)
			}
			var got httpx.HTTPError
			if err := json.UnmarshalRead(rec.Body, &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			want := httpx.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
			if got != want {
				t.Errorf("HTTPError = %+v, want %+v", got, want)
			}
		})
	}
}

func TestReadTodo_WrappedNotFound(t *testing.T) {
	notFoundErr := fmt.Errorf("todo 5: %w", model.ErrNotFound)
	c := api.New(failingStore{err: notFoundErr})
	rec := httptest.NewRecorder()
	c.ReadTodo(rec, newRequest(t, "5", ""))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNotFound, rec.Body)
	}
	var got httpx.HTTPError
	if err := json.UnmarshalRead(rec.Body, &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	want := httpx.HTTPError{
		Code:    http.StatusNotFound,
		Message: notFoundErr.Error(),
	}
	if got != want {
		t.Errorf("HTTPError = %+v, want %+v", got, want)
	}
}
