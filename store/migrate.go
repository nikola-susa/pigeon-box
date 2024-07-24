package store

import (
	"context"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/maragudk/migrate"
	"io/fs"
)

//go:embed sqlite/*.sql
var sqlite embed.FS

func (s *Store) Migrate() error {

	opt, err := setup(s.db)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	err = migrate.New(opt).MigrateUp(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func setup(db *sqlx.DB) (migrate.Options, error) {

	opt := migrate.Options{
		DB:    db.DB,
		FS:    sqlite,
		Table: "migrations",
	}

	switch db.DriverName() {
	case "libsql":
		folder, _ := fs.Sub(sqlite, "sqlite")
		opt.FS = folder
	default:
		return migrate.Options{}, fmt.Errorf("unsupported database driver: %s", db.DriverName())
	}

	return opt, nil
}
