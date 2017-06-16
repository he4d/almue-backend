package store

import (
	"database/sql"
	"log"

	"github.com/he4d/almue/model"
)

// Store must be implemented by all data stores
type Store interface {
	GetFloor(floorID int64) (*model.Floor, error)

	GetFloorList() ([]*model.Floor, error)

	CreateFloor(*model.Floor) (int64, error)

	UpdateFloor(*model.Floor) error

	DeleteFloor(floorID int64) error

	GetShutterByFloor(shutterID int64, floorID int64) (*model.Shutter, error)

	GetShutterListOfFloor(int64) ([]*model.Shutter, error)

	GetAllShutters() ([]*model.Shutter, error)

	CreateShutter(*model.Shutter) (int64, error)

	UpdateShutter(*model.Shutter) error

	DeleteShutter(shutterID int64) error

	GetLightingByFloor(lightingID int64, floorID int64) (*model.Lighting, error)

	GetLightingListOfFloor(floorID int64) ([]*model.Lighting, error)

	GetAllLightings() ([]*model.Lighting, error)

	CreateLighting(*model.Lighting) (int64, error)

	UpdateLighting(*model.Lighting) error

	DeleteLighting(int64) error
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
