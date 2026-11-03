package httputil

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Decode reads the request body as JSON into v, rejecting unknown fields.
func Decode(r *http.Request, v any) error {
	return json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(true))
}

// JSON writes v as a json response
func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, v); err != nil {
		slog.ErrorContext(r.Context(), "writing JSON response", "error", err)
	}
}

// NewError writes err as a JSON error response, including the request ID
func NewError(w http.ResponseWriter, r *http.Request, status int, err error) {
	JSON(w, r, status, HTTPError{
		Code:      status,
		Message:   err.Error(),
		RequestID: middleware.GetReqID(r.Context()),
	})
}

// HTTPError example
type HTTPError struct {
	Code      int    `json:"code" example:"400" validate:"required"`
	Message   string `json:"message" example:"status bad request" validate:"required"`
	RequestID string `json:"request_id,omitempty" example:"host/abc123-000001"`
}
