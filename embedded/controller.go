package embedded

import (
	"log"

	"periph.io/x/periph/host"

	"github.com/he4d/almue/model"
	"github.com/he4d/simplejack"
)

//StateSyncChannels holds the channels for the synchronization of the state
type StateSyncChannels struct {
	State chan string
	Quit  chan bool
}

//DeviceController must be implemented by every controller that should control the devices
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

	GetShutterStateSyncChannels(shutterID int64) (*StateSyncChannels, error)

	GetLightingStateSyncChannels(lightingID int64) (*StateSyncChannels, error)
}

type deviceController struct {
	shutters  map[int64]*shutter
	lightings map[int64]*lighting
	simulate  bool
	logger    *simplejack.Logger
}

//New creates a new DeviceController and returns it
//if true is passed to the simulate argument it runs without gpio acces
func New(simulate bool, logger *simplejack.Logger) DeviceController {
	if !simulate {
		if _, err := host.Init(); err != nil {
			log.Fatal(err)
		}
	}
	return &deviceController{
		shutters:  make(map[int64]*shutter),
		lightings: make(map[int64]*lighting),
		simulate:  simulate,
		logger:    logger,
	}
}
