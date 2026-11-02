package model

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

// Todo example
type Todo struct {
	ID    int    `json:"id" example:"1" format:"int64" validate:"required"`
	Title string `json:"title" example:"buy milk" validate:"required"`
	Done  bool   `json:"done" example:"false" validate:"required"`
}

// example
var (
	ErrTitleInvalid = errors.New("title is empty")
)

// AddTodo example
type AddTodo struct {
	Title string `json:"title" example:"buy milk" validate:"required"`
}

// Validation example
func (a AddTodo) Validation() error {
	if len(a.Title) == 0 {
		return ErrTitleInvalid
	}
	return nil
}

// UpdateTodo example, nil fields are left unchanged
type UpdateTodo struct {
	Title *string `json:"title,omitempty" example:"buy milk"`
	Done  *bool   `json:"done,omitempty" example:"true"`
}

// Validation example
func (u UpdateTodo) Validation() error {
	if u.Title != nil && len(*u.Title) == 0 {
		return ErrTitleInvalid
	}
	return nil
}

// TodosAll example
func TodosAll() []Todo {
	mu.Lock()
	defer mu.Unlock()
	return slices.Clone(todos)
}

// TodoOne example
func TodoOne(id int) (Todo, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, v := range todos {
		if id == v.ID {
			return v, nil
		}
	}
	return Todo{}, notFound(id)
}

// Insert example
func Insert(title string) Todo {
	mu.Lock()
	defer mu.Unlock()
	todoMaxID++
	todo := Todo{ID: todoMaxID, Title: title}
	todos = append(todos, todo)
	return todo
}

// Update example
func Update(id int, u UpdateTodo) (Todo, error) {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range todos {
		if id == v.ID {
			if u.Title != nil {
				todos[k].Title = *u.Title
			}
			if u.Done != nil {
				todos[k].Done = *u.Done
			}
			return todos[k], nil
		}
	}
	return Todo{}, notFound(id)
}

// Delete example
func Delete(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range todos {
		if id == v.ID {
			todos = slices.Delete(todos, k, k+1)
			return nil
		}
	}
	return notFound(id)
}

func notFound(id int) error {
	return fmt.Errorf("todo id=%d is not found", id)
}

var mu sync.Mutex
var todoMaxID = 2
var todos = []Todo{
	{ID: 1, Title: "buy milk"},
	{ID: 2, Title: "walk the dog", Done: true},
}
