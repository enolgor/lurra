package httpx

import (
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/http"
)

// Decode reads the request body as JSON into v, rejecting unknown fields.
// Errors are a StatusError: 413 for oversized bodies, 400 otherwise.
func Decode(r *http.Request, v any) error {
	err := json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(true))
	if err == nil {
		return nil
	}
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return WithStatus(http.StatusRequestEntityTooLarge, err)
	}
	return WithStatus(http.StatusBadRequest, err)
}

// JSON writes v as a json response
func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, v); err != nil { // coverage-ignore: only fails when the client disconnects mid-write, or for types handlers never pass
		slog.ErrorContext(r.Context(), "writing JSON response", "error", err)
	}
}
