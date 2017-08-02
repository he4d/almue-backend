package embedded

import (
	"periph.io/x/periph/host"

	"sync"

	"github.com/he4d/simplejack"
)

// Controller holds all necessary fields for the Controller
type Controller struct {
	shuttersLock  sync.RWMutex
	shutters      map[int64]*shutter
	lightingsLock sync.RWMutex
	lightings     map[int64]*lighting
	simulate      bool
	logger        *simplejack.Logger
	stateStore    DeviceStateStore
}

//New creates a new DeviceController and returns it
//if true is passed to the simulate argument it runs without gpio acces
func New(logger *simplejack.Logger, stateStore DeviceStateStore, simulate bool) (*Controller, error) {
	if !simulate {
		if _, err := host.Init(); err != nil {
			return nil, err
		}
	}

	controller := &Controller{
		shutters:   make(map[int64]*shutter),
		lightings:  make(map[int64]*lighting),
		simulate:   simulate,
		stateStore: stateStore,
		logger:     logger,
	}

	return controller, nil
}
