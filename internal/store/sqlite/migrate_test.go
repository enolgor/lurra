package sqlite

import (
	"context"
	"database/sql"
	"testing"
)

// TestMigrationsUpDown applies all migrations, rolls them all back and applies them again,
// so every Down section is known to work.
func TestMigrationsUpDown(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1) // every connection to :memory: is a separate database

	provider, err := newProvider(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := provider.Up(ctx); err != nil {
		t.Fatalf("up: %v", err)
	}
	if _, err := provider.DownTo(ctx, 0); err != nil {
		t.Fatalf("down: %v", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		t.Fatalf("up again: %v", err)
	}
}

func TestMigrate_NilDB(t *testing.T) {
	if err := migrate(t.Context(), nil); err == nil {
		t.Fatal("migrate(nil) error = nil, want an error")
	}
}
