package controller

import (
	"net/http"

	"github.com/enolgor/lurra/internal/httputil"
	"github.com/enolgor/lurra/internal/model"
)

// ListTodos godoc
//
//	@Summary		List todos
//	@Description	get all todos
//	@Tags			todos
//	@Produce		json
//	@Success		200	{array}		model.Todo
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/todos [get]
func (c *Controller) ListTodos(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, r, http.StatusOK, model.TodosAll())
}

// ShowTodo godoc
//
//	@Summary		Show a todo
//	@Description	get todo by ID
//	@Tags			todos
//	@Produce		json
//	@Param			id	path		int	true	"Todo ID"
//	@Success		200	{object}	model.Todo
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/todos/{id} [get]
func (c *Controller) ShowTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	todo, err := model.TodoOne(id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, r, http.StatusOK, todo)
}

// AddTodo godoc
//
//	@Summary		Add a todo
//	@Description	add by json todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			todo	body		model.AddTodo	true	"Add todo"
//	@Success		201		{object}	model.Todo
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		413		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/todos [post]
func (c *Controller) AddTodo(w http.ResponseWriter, r *http.Request) {
	var addTodo model.AddTodo
	if err := httputil.Decode(r, &addTodo); err != nil {
		writeError(w, r, err)
		return
	}
	if err := addTodo.Validation(); err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, r, http.StatusCreated, model.Insert(addTodo.Title))
}

// UpdateTodo godoc
//
//	@Summary		Update a todo
//	@Description	update title and done state by json todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Todo ID"
//	@Param			todo	body		model.UpdateTodo	true	"Update todo"
//	@Success		200		{object}	model.Todo
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		413		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/todos/{id} [patch]
func (c *Controller) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var updateTodo model.UpdateTodo
	if err := httputil.Decode(r, &updateTodo); err != nil {
		writeError(w, r, err)
		return
	}
	if err := updateTodo.Validation(); err != nil {
		writeError(w, r, err)
		return
	}
	todo, err := model.Update(id, updateTodo)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, r, http.StatusOK, todo)
}

// DeleteTodo godoc
//
//	@Summary		Delete a todo
//	@Description	delete by todo ID
//	@Tags			todos
//	@Produce		json
//	@Param			id	path	int	true	"Todo ID"
//	@Success		204
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/todos/{id} [delete]
func (c *Controller) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if err := model.Delete(id); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
