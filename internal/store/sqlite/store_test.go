package sqlite_test

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/enolgor/lurra/internal/model"
	"github.com/enolgor/lurra/internal/store/sqlite"
)

func TestOpen_Errors(t *testing.T) {
	t.Run("path is a directory", func(t *testing.T) {
		_, err := sqlite.Open(t.Context(), t.TempDir())
		if err == nil || !strings.Contains(err.Error(), "opening database") {
			t.Fatalf("Open() error = %v, want an opening database error", err)
		}
	})
	t.Run("conflicting table", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "app.db")
		db, err := sql.Open("sqlite3", "file:"+path)
		if err != nil {
			t.Fatalf("open db: %v", err)
		}
		if _, err = db.ExecContext(t.Context(), `CREATE TABLE todos (x INTEGER)`); err != nil {
			t.Fatalf("create todos table: %v", err)
		}
		if err = db.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
		_, err = sqlite.Open(t.Context(), path)
		if err == nil || !strings.Contains(err.Error(), "migrating database") {
			t.Fatalf("Open() error = %v, want a migration error", err)
		}
	})
}

// TestStore_ClosedDatabase checks that a real database error is not mistaken for
// model.ErrNotFound, which would turn a 500 into a misleading 404.
func TestStore_ClosedDatabase(t *testing.T) {
	store, err := sqlite.Open(t.Context(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	tests := []struct {
		name string
		call func(s *sqlite.Store) error
	}{
		{name: "read", call: func(s *sqlite.Store) error {
			_, err := s.ReadTodo(t.Context(), 1)
			return err
		}},
		{name: "update", call: func(s *sqlite.Store) error {
			_, err := s.UpdateTodo(t.Context(), 1, model.UpdateTodo{})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(store)
			if err == nil {
				t.Fatal("error = nil, want an error")
			}
			if errors.Is(err, model.ErrNotFound) {
				t.Errorf("error = %v, want an error that is not ErrNotFound", err)
			}
		})
	}
}

func newTestStore(t *testing.T) *sqlite.Store {
	t.Helper()
	s, err := sqlite.Open(t.Context(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestStore_Concurrent(t *testing.T) {
	store := newTestStore(t)
	const n = 50

	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			if _, err := store.CreateTodo(t.Context(), "buy milk"); err != nil {
				t.Errorf("create (concurrent): %v", err)
			}
		})
	}
	wg.Wait()

	todos, err := store.ListTodos(t.Context())
	if err != nil {
		t.Fatalf("todos: %v", err)
	}
	if len(todos) != n {
		t.Fatalf("len(todos) = %d, want %d", len(todos), n)
	}
	seen := make(map[int]bool)
	for _, todo := range todos {
		seen[todo.ID] = true
	}
	if len(seen) != n {
		t.Fatalf("len(seen) = %d, want %d", len(seen), n)
	}
}

func TestStore_CRUD(t *testing.T) {
	store := newTestStore(t)
	ctx := t.Context()

	created, err := store.CreateTodo(ctx, "buy milk")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if want := (model.Todo{ID: created.ID, Title: "buy milk"}); created != want || created.ID == 0 {
		t.Fatalf("created = %+v, want %+v with a non-zero ID", created, want)
	}

	read, err := store.ReadTodo(ctx, created.ID)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if read != created {
		t.Errorf("read = %+v, want %+v", read, created)
	}

	// Partial updates: each one changes only the field that is set.
	done := true
	updated, err := store.UpdateTodo(ctx, created.ID, model.UpdateTodo{Done: &done})
	if err != nil {
		t.Fatalf("update done: %v", err)
	}
	if want := (model.Todo{ID: created.ID, Title: "buy milk", Done: true}); updated != want {
		t.Errorf("after done update = %+v, want %+v", updated, want)
	}
	title := "buy oat milk"
	updated, err = store.UpdateTodo(ctx, created.ID, model.UpdateTodo{Title: &title})
	if err != nil {
		t.Fatalf("update title: %v", err)
	}
	if want := (model.Todo{ID: created.ID, Title: "buy oat milk", Done: true}); updated != want {
		t.Errorf("after title update = %+v, want %+v", updated, want)
	}

	if err = store.DeleteTodo(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	todos, err := store.ListTodos(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("todos after delete = %+v, want none", todos)
	}
}

func TestStore_NotFound(t *testing.T) {
	store := newTestStore(t)
	const id = 999

	tests := []struct {
		name string
		call func() error
	}{
		{name: "read", call: func() error {
			_, err := store.ReadTodo(t.Context(), id)
			return err
		}},
		{name: "update", call: func() error {
			_, err := store.UpdateTodo(t.Context(), id, model.UpdateTodo{})
			return err
		}},
		{name: "delete", call: func() error {
			return store.DeleteTodo(t.Context(), id)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if !errors.Is(err, model.ErrNotFound) {
				t.Fatalf("error = %v, want model.ErrNotFound", err)
			}
			if want := "todo id=999 is not found"; err.Error() != want {
				t.Errorf("error = %q, want %q", err, want)
			}
		})
	}
}
