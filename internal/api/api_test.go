package api_test

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enolgor/lurra/internal/api"
	"github.com/enolgor/lurra/internal/httpx"
	"github.com/enolgor/lurra/internal/model"
)

// TestController_InvalidID checks that a non-numeric {id} is rejected with 400 before the
// store is called; the zero failingStore would panic if any of its methods ran.
func TestController_InvalidID(t *testing.T) {
	c := api.New(failingStore{})

	tests := []struct {
		name    string
		handler http.HandlerFunc
		body    string
	}{
		{name: "ReadTodo", handler: c.ReadTodo},
		{name: "UpdateTodo", handler: c.UpdateTodo, body: `{"done":true}`},
		{name: "DeleteTodo", handler: c.DeleteTodo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.handler(rec, newRequest(t, "abc", tt.body))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body)
			}
			var got httpx.HTTPError
			if err := json.UnmarshalRead(rec.Body, &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got.Code != http.StatusBadRequest || got.Message == "" {
				t.Errorf("HTTPError = %+v, want code %d and a message", got, http.StatusBadRequest)
			}
		})
	}
}

// TestController_InvalidTitle checks that model.ErrTitleInvalid from validation maps to 400
// with the validation message.
func TestController_InvalidTitle(t *testing.T) {
	c := api.New(failingStore{})

	tests := []struct {
		name    string
		handler http.HandlerFunc
		id      string
	}{
		{name: "CreateTodo", handler: c.CreateTodo},
		{name: "UpdateTodo", handler: c.UpdateTodo, id: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.handler(rec, newRequest(t, tt.id, `{"title":""}`))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body)
			}
			var got httpx.HTTPError
			if err := json.UnmarshalRead(rec.Body, &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			want := httpx.HTTPError{Code: http.StatusBadRequest, Message: model.ErrTitleInvalid.Error()}
			if got != want {
				t.Errorf("HTTPError = %+v, want %+v", got, want)
			}
		})
	}
}
