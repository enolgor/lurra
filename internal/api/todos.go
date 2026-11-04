package api

import (
	"net/http"

	"github.com/enolgor/lurra/internal/httpx"
	"github.com/enolgor/lurra/internal/model"
)

// @Summary		List todos
// @Produce		json
// @Success		200	{array}		model.Todo
// @Failure		500	{object}	httpx.HTTPError
// @Router			/todos [get]
func (c *Controller) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := c.todos.ListTodos(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, todos)
}

// @Summary		Show a todo
// @Produce		json
// @Param			id	path		int	true	"Todo ID"
// @Success		200	{object}	model.Todo
// @Failure		400	{object}	httpx.HTTPError
// @Failure		404	{object}	httpx.HTTPError
// @Failure		500	{object}	httpx.HTTPError
// @Router			/todos/{id} [get]
func (c *Controller) ShowTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	todo, err := c.todos.GetTodo(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, todo)
}

// @Summary		Add a todo
// @Accept			json
// @Produce		json
// @Param			todo	body		model.AddTodo	true	"Add todo"
// @Success		201		{object}	model.Todo
// @Failure		400		{object}	httpx.HTTPError
// @Failure		413		{object}	httpx.HTTPError
// @Failure		500		{object}	httpx.HTTPError
// @Router			/todos [post]
func (c *Controller) AddTodo(w http.ResponseWriter, r *http.Request) {
	var addTodo model.AddTodo
	if err := httpx.Decode(r, &addTodo); err != nil {
		writeError(w, r, err)
		return
	}
	if err := addTodo.Validation(); err != nil {
		writeError(w, r, err)
		return
	}
	todo, err := c.todos.CreateTodo(r.Context(), addTodo.Title)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, todo)
}

// @Summary		Update a todo
// @Description	Only the fields present in the body are changed.
// @Accept			json
// @Produce		json
// @Param			id		path		int					true	"Todo ID"
// @Param			todo	body		model.UpdateTodo	true	"Update todo"
// @Success		200		{object}	model.Todo
// @Failure		400		{object}	httpx.HTTPError
// @Failure		404		{object}	httpx.HTTPError
// @Failure		413		{object}	httpx.HTTPError
// @Failure		500		{object}	httpx.HTTPError
// @Router			/todos/{id} [patch]
func (c *Controller) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var updateTodo model.UpdateTodo
	if err := httpx.Decode(r, &updateTodo); err != nil {
		writeError(w, r, err)
		return
	}
	if err := updateTodo.Validation(); err != nil {
		writeError(w, r, err)
		return
	}
	todo, err := c.todos.UpdateTodo(r.Context(), id, updateTodo)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, todo)
}

// @Summary		Delete a todo
// @Produce		json
// @Param			id	path	int	true	"Todo ID"
// @Success		204
// @Failure		400	{object}	httpx.HTTPError
// @Failure		404	{object}	httpx.HTTPError
// @Failure		500	{object}	httpx.HTTPError
// @Router			/todos/{id} [delete]
func (c *Controller) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if err := c.todos.DeleteTodo(r.Context(), id); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
