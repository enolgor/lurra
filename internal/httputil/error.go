package httputil

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

// StatusError is an error with the HTTP status to respond with.
type StatusError struct {
	Status int
	Err    error
}

func (e *StatusError) Error() string { return e.Err.Error() }
func (e *StatusError) Unwrap() error { return e.Err }

// WithStatus wraps err with the HTTP status to respond with.
func WithStatus(status int, err error) error {
	return &StatusError{Status: status, Err: err}
}

// NewError writes err as a JSON error response, including the request ID, and adds err to the request log.
// Errors without a StatusError are sent as a generic 500, so internal details don't leak.
func NewError(w http.ResponseWriter, r *http.Request, err error) {
	_ = httplog.SetError(r.Context(), err)
	status, msg := http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	if se, ok := errors.AsType[*StatusError](err); ok {
		status, msg = se.Status, se.Error()
	}
	JSON(w, r, status, HTTPError{
		Code:      status,
		Message:   msg,
		RequestID: middleware.GetReqID(r.Context()),
	})
}

// HTTPError is the JSON body of error responses.
type HTTPError struct {
	Code      int    `json:"code" example:"400" validate:"required"`
	Message   string `json:"message" example:"status bad request" validate:"required"`
	RequestID string `json:"request_id,omitempty" example:"host/abc123-000001"`
}
