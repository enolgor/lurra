package httpx

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

// StatusError is an error with the HTTP status to respond with and, when the status alone
// doesn't say which error it is, a machine-readable error code for clients to act on.
type StatusError struct {
	Status    int
	ErrorCode string
	Err       error
}

func (e *StatusError) Error() string { return e.Err.Error() }
func (e *StatusError) Unwrap() error { return e.Err }

// WithStatus wraps err with the HTTP status to respond with.
func WithStatus(status int, err error) error {
	return &StatusError{Status: status, Err: err}
}

// WithErrorCode wraps err with the HTTP status and an error code, for errors that share a status
// with others (e.g. two kinds of 400), so clients can tell them apart. Add every new code to the
// enums tag of HTTPError.ErrorCode, so the generated frontend type lists it.
func WithErrorCode(status int, code string, err error) error {
	return &StatusError{Status: status, ErrorCode: code, Err: err}
}

// NewError writes err as a JSON error response, including the request ID, and adds err to the request log.
// Errors without a StatusError are sent as a generic 500, so internal details don't leak.
func NewError(w http.ResponseWriter, r *http.Request, err error) {
	_ = httplog.SetError(r.Context(), err)
	body := HTTPError{
		Code:      http.StatusInternalServerError,
		Message:   http.StatusText(http.StatusInternalServerError),
		RequestID: middleware.GetReqID(r.Context()),
	}
	if se, ok := errors.AsType[*StatusError](err); ok {
		body.Code, body.Message, body.ErrorCode = se.Status, se.Error(), se.ErrorCode
	}
	JSON(w, r, body.Code, body)
}

// HTTPError is the JSON body of error responses.
type HTTPError struct {
	Code      int    `json:"code" example:"400" validate:"required"`
	Message   string `json:"message" example:"title is empty" validate:"required"`
	RequestID string `json:"request_id,omitempty" example:"host/abc123-000001"`
	// Identifies the error when code alone doesn't (e.g. two kinds of 400).
	ErrorCode string `json:"error,omitempty" enums:"title_invalid" example:"title_invalid"`
}
