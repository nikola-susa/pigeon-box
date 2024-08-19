package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/nikola-susa/pigeon-box/config"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type Store struct {
	db     *sqlx.DB
	config *config.Config
	ctx    context.Context
}

func New(config *config.Config) (*Store, error) {

	dbx, err := sqlx.Connect("libsql", config.Database.URL)
	if err != nil {
		return nil, err
	}

	_, err = dbx.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return &Store{
		db:     dbx,
		config: config,
	}, nil
}

func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) Close() error {
	return s.db.Close()
}
