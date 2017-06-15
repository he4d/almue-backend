package store

import (
	"database/sql"
	"log"

	"github.com/he4d/almue/model"
)

// Store must be implemented by all data stores
type Store interface {
	GetFloor(int64) (*model.Floor, error)

	GetFloorList() ([]*model.Floor, error)

	CreateFloor(*model.Floor) error

	UpdateFloor(*model.Floor) error

	DeleteFloor(*model.Floor) error

	GetShutterByFloor(int64, int64) (*model.Shutter, error)

	GetShutterListOfFloor(int64) ([]*model.Shutter, error)

	GetAllShutters() ([]*model.Shutter, error)

	CreateShutter(*model.Shutter) error

	UpdateShutter(*model.Shutter) error

	DeleteShutter(*model.Shutter) error

	GetLightingByFloor(int64, int64) (*model.Lighting, error)

	GetLightingListOfFloor(int64) ([]*model.Lighting, error)

	GetAllLightings() ([]*model.Lighting, error)

	CreateLighting(*model.Lighting) error

	UpdateLighting(*model.Lighting) error

	DeleteLighting(*model.Lighting) error
}

type datastore struct {
	*sql.DB
}

// New returns a new datastore that is completely initialized
func New(path string) Store {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Println(err)
		log.Fatalln("database connection failed")
	}

	if err := setupDatabase(db); err != nil {
		log.Println(err)
		log.Fatalln("migration failed")
	}
	return &datastore{DB: db}
}

func setupDatabase(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return Migrate(db)
}
