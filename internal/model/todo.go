package model

import (
	"context"
	"errors"
)

// TodoStore persists todos. Methods return ErrNotFound (wrapped) for unknown IDs.
// Implementations live in internal/store, e.g. sqlite.Store.
type TodoStore interface {
	ListTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	CreateTodo(ctx context.Context, title string) (Todo, error)
	UpdateTodo(ctx context.Context, id int, u UpdateTodo) (Todo, error)
	DeleteTodo(ctx context.Context, id int) error
}

// Todo is a stored todo item.
type Todo struct {
	ID    int    `json:"id" example:"1" format:"int64" validate:"required"`
	Title string `json:"title" example:"buy milk" validate:"required"`
	Done  bool   `json:"done" example:"false" validate:"required"`
}

var (
	ErrNotFound     = errors.New("not found")
	ErrTitleInvalid = errors.New("title is empty")
)

// AddTodo is the body for creating a todo.
type AddTodo struct {
	Title string `json:"title" example:"buy milk" validate:"required"`
}

// Validation checks the fields of a new todo.
func (a AddTodo) Validation() error {
	if len(a.Title) == 0 {
		return ErrTitleInvalid
	}
	return nil
}

// UpdateTodo is the body for updating a todo; nil fields are left unchanged.
type UpdateTodo struct {
	Title *string `json:"title,omitempty" example:"buy milk"`
	Done  *bool   `json:"done,omitempty" example:"true"`
}

// Validation checks the fields that are set.
func (u UpdateTodo) Validation() error {
	if u.Title != nil && len(*u.Title) == 0 {
		return ErrTitleInvalid
	}
	return nil
}
