// Package sqlite is a SQLite store using database/sql. The schema is managed by embedded migrations.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/ncruces/go-sqlite3/driver" // registers the "sqlite3" database/sql driver
)

// Store keeps data in a SQLite database and is safe for concurrent use.
type Store struct {
	db *sql.DB
}

// pragmas are set on every new connection; busy_timeout must come first.
var pragmas = []string{
	"busy_timeout(5000)",
	"journal_mode(wal)",
	"synchronous(normal)",
	"foreign_keys(on)",
	"temp_store(memory)", // sorts/temp tables in RAM: faster, and no writable temp dir needed (read-only root fs)
}

// Open opens (or creates) the database at path and applies pending migrations.
func Open(ctx context.Context, path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	// _txlock=immediate takes the write lock when a transaction starts, avoiding
	// "database is locked" errors when a read transaction later tries to write.
	query := url.Values{"_pragma": pragmas, "_txlock": {"immediate"}}
	db, err := sql.Open("sqlite3", "file:"+path+"?"+query.Encode())
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrating database: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}
