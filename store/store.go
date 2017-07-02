package store

import (
	"database/sql"

	"github.com/he4d/simplejack"
)

type Datastore struct {
	*sql.DB
	logger *simplejack.Logger
}

// New returns a new datastore that is completely initialized
func New(path string, logger *simplejack.Logger) (*Datastore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := setupDatabase(db); err != nil {
		return nil, err
	}
	return &Datastore{DB: db}, nil
}

func setupDatabase(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return Migrate(db)
}
