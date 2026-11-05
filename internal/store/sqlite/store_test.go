package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/enolgor/lurra/internal/model"
	"github.com/enolgor/lurra/internal/store/sqlite"
)

func TestOpen_Errors(t *testing.T) {
	t.Run("parent is a file", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "file")
		if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
			t.Fatalf("write file: %v", err)
		}
		_, err := sqlite.Open(t.Context(), filepath.Join(file, "app.db"))
		if err == nil || !strings.Contains(err.Error(), "creating database directory") {
			t.Fatalf("Open() error = %v, want a directory error", err)
		}
	})
	t.Run("path is a directory", func(t *testing.T) {
		_, err := sqlite.Open(t.Context(), t.TempDir())
		if err == nil || !strings.Contains(err.Error(), "opening database") {
			t.Fatalf("Open() error = %v, want an opening database error", err)
		}
	})
	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if _, err := sqlite.Open(ctx, filepath.Join(t.TempDir(), "app.db")); !errors.Is(err, context.Canceled) {
			t.Fatalf("Open() error = %v, want a context canceled error", err)
		}
	})
	t.Run("conflicting table", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "app.db")
		db, err := sql.Open("sqlite3", "file:"+path)
		if err != nil {
			t.Fatalf("open db: %v", err)
		}
		if _, err = db.Exec(`CREATE TABLE todos (x INTEGER)`); err != nil {
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
		{name: "list", call: func(s *sqlite.Store) error {
			_, err := s.ListTodos(t.Context())
			return err
		}},
		{name: "create", call: func(s *sqlite.Store) error {
			_, err := s.CreateTodo(t.Context(), "buy milk")
			return err
		}},
		{name: "read", call: func(s *sqlite.Store) error {
			_, err := s.ReadTodo(t.Context(), 1)
			return err
		}},
		{name: "update", call: func(s *sqlite.Store) error {
			_, err := s.UpdateTodo(t.Context(), 1, model.UpdateTodo{})
			return err
		}},
		{name: "delete", call: func(s *sqlite.Store) error {
			return s.DeleteTodo(t.Context(), 1)
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
