package almue

import "github.com/he4d/almue/model"

// DeviceStore must be implemented by the device data store
type DeviceStore interface {
	GetFloor(floorID int64) (*model.Floor, error)

	GetFloorList() ([]*model.Floor, error)

	CreateFloor(*model.Floor) (int64, error)

	UpdateFloor(*model.Floor) error

	DeleteFloor(floorID int64) error

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

	GetBackup() ([]byte, error)

	RestoreBackup([]byte) error
}

// DeviceController must be implemented by the device controller
type DeviceController interface {
	RegisterShutters(shutters ...*model.Shutter) error

	UnregisterShutter(shutterID int64) error

	UpdateShutter(diffs model.DifferenceType, updatedShutter *model.Shutter) error

	RegisterLightings(lightings ...*model.Lighting) error

	UnregisterLighting(lightingID int64) error

	UpdateLighting(diffs model.DifferenceType, updatedLighting *model.Lighting) error

	OpenShutter(shutterID int64) error

	CloseShutter(shutterID int64) error

	StopShutter(shutterID int64) error

	TurnLightingOn(lightingID int64) error

	TurnLightingOff(lightingID int64) error

	ScheduleShutterJobs(shutter *model.Shutter) error

	UnscheduleShutterJobs(shutterID int64) error

	ScheduleLightingJobs(lighting *model.Lighting) error

	UnscheduleLightingJobs(lightingID int64) error
}
