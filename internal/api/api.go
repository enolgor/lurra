package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enolgor/lurra/internal/httpx"
	"github.com/enolgor/lurra/internal/model"
	"github.com/go-chi/chi/v5"
)

// General API info for swag (make bindings reads it from this file).

//	@title			Todo API
//	@version		1.0

// Controller holds the HTTP handlers and their dependencies.
type Controller struct {
	todos model.TodoStore
}

// New returns a controller using todos for storage.
func New(todos model.TodoStore) *Controller {
	return &Controller{todos: todos}
}

// writeError maps domain errors to HTTP statuses and writes the error response. An error that
// shares its status with others also gets an error code (listed in httpx.HTTPError's enums).
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		err = httpx.WithStatus(http.StatusNotFound, err)
	case errors.Is(err, model.ErrTitleInvalid):
		err = httpx.WithErrorCode(http.StatusBadRequest, "title_invalid", err)
	}
	httpx.NewError(w, r, err)
}

// pathID parses the {id} path parameter.
func pathID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, httpx.WithStatus(http.StatusBadRequest, err)
	}
	return id, nil
}
