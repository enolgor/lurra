package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/enolgor/lurra/internal/model"
)

var _ model.TodoStore = (*Store)(nil)

func (s *Store) ListTodos(ctx context.Context) ([]model.Todo, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, title, done FROM todos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil { // coverage-ignore
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *Store) CreateTodo(ctx context.Context, title string) (model.Todo, error) {
	row := s.db.QueryRowContext(ctx, `INSERT INTO todos (title) VALUES (?) RETURNING id, title, done`, title)
	return scanTodo(row, 0)
}

func (s *Store) ReadTodo(ctx context.Context, id int) (model.Todo, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, title, done FROM todos WHERE id = ?`, id)
	return scanTodo(row, id)
}

// UpdateTodo changes only the fields that are set; nil pointers bind as NULL and COALESCE keeps the old value.
func (s *Store) UpdateTodo(ctx context.Context, id int, u model.UpdateTodo) (model.Todo, error) {
	row := s.db.QueryRowContext(ctx, `
		UPDATE todos SET title = COALESCE(?, title), done = COALESCE(?, done)
		WHERE id = ?
		RETURNING id, title, done`, u.Title, u.Done, id)
	return scanTodo(row, id)
}

func (s *Store) DeleteTodo(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil { // coverage-ignore
		return err
	}
	if n == 0 {
		return notFound(id)
	}
	return nil
}

// scanTodo scans a single todo row, turning sql.ErrNoRows into model.ErrNotFound for id.
func scanTodo(row *sql.Row, id int) (model.Todo, error) {
	var t model.Todo
	err := row.Scan(&t.ID, &t.Title, &t.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Todo{}, notFound(id)
	}
	return t, err
}

func notFound(id int) error {
	return fmt.Errorf("todo id=%d is %w", id, model.ErrNotFound)
}
