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
	return &Datastore{DB: db, logger: logger}, nil
}

func (s *Datastore) Shutdown() {
	s.logger.Info.Print("Shutting down the datastore...")
	if err := s.Close(); err != nil {
		s.logger.Error.Printf("Could not shutdown the store: %v", err)
		return
	}
}

func setupDatabase(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return Migrate(db)
}
