package embedded

import (
	"periph.io/x/periph/host"

	"github.com/he4d/simplejack"
)

type EmbeddedController struct {
	shutters   map[int64]*shutter
	lightings  map[int64]*lighting
	simulate   bool
	logger     *simplejack.Logger
	stateStore DeviceStateStore
}

//New creates a new DeviceController and returns it
//if true is passed to the simulate argument it runs without gpio acces
func New(logger *simplejack.Logger, stateStore DeviceStateStore, simulate bool) (*EmbeddedController, error) {
	if !simulate {
		if _, err := host.Init(); err != nil {
			return nil, err
		}
	}

	controller := &EmbeddedController{
		shutters:   make(map[int64]*shutter),
		lightings:  make(map[int64]*lighting),
		simulate:   simulate,
		stateStore: stateStore,
		logger:     logger,
	}

	return controller, nil
}
