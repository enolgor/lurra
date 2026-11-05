package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// Migrations are goose SQL files named NNNNN_description.sql, each with "-- +goose Up" and
// "-- +goose Down" sections. Goose applies them in order, each in its own transaction, and
// records applied versions in goose_db_version. Never edit an applied migration; add a new one.
//
//go:embed migrations/*.sql
var migrations embed.FS

func migrate(ctx context.Context, db *sql.DB) error {
	provider, err := newProvider(db)
	if err != nil {
		return err
	}
	_, err = provider.Up(ctx)
	return err
}

func newProvider(db *sql.DB) (*goose.Provider, error) {
	fsys, err := fs.Sub(migrations, "migrations")
	if err != nil { // coverage-ignore
		return nil, err
	}
	return goose.NewProvider(goose.DialectSQLite3, db, fsys, goose.WithSlog(slog.Default()))
}
