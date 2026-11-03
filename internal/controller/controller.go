package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enolgor/lurra/internal/httputil"
	"github.com/enolgor/lurra/internal/model"
	"github.com/go-chi/chi/v5"
)

// Controller example
type Controller struct{}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// writeError maps domain errors to HTTP statuses and writes the error response.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		err = httputil.WithStatus(http.StatusNotFound, err)
	case errors.Is(err, model.ErrTitleInvalid):
		err = httputil.WithStatus(http.StatusBadRequest, err)
	}
	httputil.NewError(w, r, err)
}

// pathID parses the {id} path parameter.
func pathID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, httputil.WithStatus(http.StatusBadRequest, err)
	}
	return id, nil
}
