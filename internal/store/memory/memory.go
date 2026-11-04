// Package memory is an in-memory store, seeded with example data. Its data is lost on restart.
package memory

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/enolgor/lurra/internal/model"
)

var _ model.TodoStore = (*Store)(nil)

// Store keeps all data in memory and is safe for concurrent use.
type Store struct {
	mu    sync.Mutex
	maxID int
	todos []model.Todo
}

// New returns a store with two example todos.
func New() *Store {
	return &Store{
		maxID: 2,
		todos: []model.Todo{
			{ID: 1, Title: "buy milk"},
			{ID: 2, Title: "walk the dog", Done: true},
		},
	}
}

func (s *Store) ListTodos(_ context.Context) ([]model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.todos), nil
}

func (s *Store) GetTodo(_ context.Context, id int) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.index(id); i >= 0 {
		return s.todos[i], nil
	}
	return model.Todo{}, notFound(id)
}

func (s *Store) CreateTodo(_ context.Context, title string) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxID++
	todo := model.Todo{ID: s.maxID, Title: title}
	s.todos = append(s.todos, todo)
	return todo, nil
}

func (s *Store) UpdateTodo(_ context.Context, id int, u model.UpdateTodo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(id)
	if i < 0 {
		return model.Todo{}, notFound(id)
	}
	if u.Title != nil {
		s.todos[i].Title = *u.Title
	}
	if u.Done != nil {
		s.todos[i].Done = *u.Done
	}
	return s.todos[i], nil
}

func (s *Store) DeleteTodo(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(id)
	if i < 0 {
		return notFound(id)
	}
	s.todos = slices.Delete(s.todos, i, i+1)
	return nil
}

// index returns the position of the todo with id, or -1. The caller must hold mu.
func (s *Store) index(id int) int {
	return slices.IndexFunc(s.todos, func(t model.Todo) bool { return t.ID == id })
}

func notFound(id int) error {
	return fmt.Errorf("todo id=%d is %w", id, model.ErrNotFound)
}
