package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// JSON writes v as a json response
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writing JSON response", "error", err)
	}
}

// NewError example
func NewError(w http.ResponseWriter, status int, err error) {
	JSON(w, status, HTTPError{
		Code:    status,
		Message: err.Error(),
	})
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400" validate:"required"`
	Message string `json:"message" example:"status bad request" validate:"required"`
}
