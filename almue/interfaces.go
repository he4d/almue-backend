package almue

import "github.com/he4d/almue/model"

// Store must be implemented by the data store
type Store interface {
	GetFloor(floorID int64) (*model.Floor, error)

	GetFloorList() ([]*model.Floor, error)

	CreateFloor(*model.Floor) (int64, error)

	UpdateFloor(*model.Floor) error

	DeleteFloor(floorID int64) error

	NumShuttersOfFloor(floorID int64) (int, error)

	NumLightingsOfFloor(floorID int64) (int, error)

	GetShutter(shutterID int64) (*model.Shutter, error)

	GetShutterList() ([]*model.Shutter, error)

	GetShutterListOfFloor(int64) ([]*model.Shutter, error)

	CreateShutter(*model.Shutter) (int64, error)

	UpdateShutter(*model.Shutter) error

	DeleteShutter(shutterID int64) error

	GetLighting(lightingID int64) (*model.Lighting, error)

	GetLightingList() ([]*model.Lighting, error)

	GetLightingListOfFloor(floorID int64) ([]*model.Lighting, error)

	CreateLighting(*model.Lighting) (int64, error)

	UpdateLighting(*model.Lighting) error

	DeleteLighting(int64) error

	UpdateLightingState(int64, string) error

	UpdateShutterState(int64, string) error
}
