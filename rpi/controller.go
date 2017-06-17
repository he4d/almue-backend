package rpi

import (
	"log"
	"time"

	"periph.io/x/periph/host"

	"fmt"

	"github.com/carlescere/scheduler"
	"github.com/he4d/almue/model"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

//StateSynchronization holds the channels for the synchronization of the state
type StateSynchronization struct {
	State chan string
	Quit  chan bool
}

//DeviceController must be implemented by every controller that should control the devices
type DeviceController interface {
	RegisterShutters(shutters ...*model.Shutter) (states map[int64]*StateSynchronization, err error)

	UnregisterShutter(shutterID int64) error

	RegisterLightings(lightings ...*model.Lighting) (states map[int64]*StateSynchronization, err error)

	UnregisterLighting(lightingID int64) error

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

type deviceController struct {
	shutters  map[int64]*shutter
	lightings map[int64]*lighting
	simulate  bool
}

//New creates a new DeviceController and returns it
//if true is passed to the simulate argument it runs without gpio acces
func New(simulate bool) DeviceController {
	if !simulate {
		if _, err := host.Init(); err != nil {
			log.Fatal(err)
		}
	}
	return &deviceController{
		shutters:  make(map[int64]*shutter),
		lightings: make(map[int64]*lighting),
		simulate:  simulate,
	}
}

func (d *deviceController) RegisterShutters(shutters ...*model.Shutter) (states map[int64]*StateSynchronization, err error) {
	states = make(map[int64]*StateSynchronization)
	for _, shutterModel := range shutters {
		var openPin gpio.PinIO
		var closePin gpio.PinIO
		if d.simulate {
			openPin = &simulatePinIO{name: shutterModel.Description, number: shutterModel.OpenPin}
			closePin = &simulatePinIO{name: shutterModel.Description, number: shutterModel.ClosePin}
		} else {
			openPin = gpioreg.ByNumber(shutterModel.OpenPin)
			closePin = gpioreg.ByNumber(shutterModel.ClosePin)
		}
		duration := time.Duration(shutterModel.CompleteWayInSeconds) * time.Second
		state := make(chan string)
		quit := make(chan bool)
		d.shutters[shutterModel.ID] = &shutter{
			ID:                  shutterModel.ID,
			openPin:             openPin,
			closePin:            closePin,
			completeWayDuration: duration,
			state:               state,
			quit:                quit,
		}
		states[shutterModel.ID] = &StateSynchronization{State: state, Quit: quit}
		if shutterModel.TimerEnabled {
			if err = d.ScheduleShutterJobs(shutterModel); err != nil {
				return nil, err
			}
		}
	}
	return states, err
}

func (d *deviceController) UnregisterShutter(shutterID int64) error {
	shutter, ok := d.shutters[shutterID]
	if !ok {
		return fmt.Errorf("Shutter with id: %d cannot be unregistered, because it doesnt exist", shutterID)
	}

	if err := d.UnscheduleShutterJobs(shutterID); err != nil {
		return err
	}

	//Stop state synchronization
	shutter.quit <- true

	delete(d.shutters, shutterID)
	return nil
}

func (d *deviceController) RegisterLightings(lightings ...*model.Lighting) (states map[int64]*StateSynchronization, err error) {
	states = make(map[int64]*StateSynchronization)
	for _, lightingModel := range lightings {
		var switchPin gpio.PinIO
		if d.simulate {
			switchPin = &simulatePinIO{name: lightingModel.Description, number: lightingModel.SwitchPin}
		} else {
			switchPin = gpioreg.ByNumber(lightingModel.SwitchPin)
		}
		state := make(chan string)
		quit := make(chan bool)
		d.lightings[lightingModel.ID] = &lighting{
			ID:        lightingModel.ID,
			switchPin: switchPin,
			state:     state,
			quit:      quit,
		}
		states[lightingModel.ID] = &StateSynchronization{State: state, Quit: quit}
		if lightingModel.TimerEnabled {
			if err = d.ScheduleLightingJobs(lightingModel); err != nil {
				return nil, err
			}
		}
	}
	return states, err
}

func (d *deviceController) UnregisterLighting(lightingID int64) error {
	lighting, ok := d.lightings[lightingID]
	if !ok {
		return fmt.Errorf("Lighting with id: %d cannot be unregistered, because it doesnt exist", lightingID)
	}
	if err := d.UnscheduleLightingJobs(lightingID); err != nil {
		return err
	}

	//Stop state synchronization
	lighting.quit <- true

	delete(d.lightings, lightingID)
	return nil
}

func (d *deviceController) OpenShutter(shutterID int64) error {
	device, err := d.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	if device.timer != nil {
		device.timer.Stop()
	}
	if err := device.closePin.Out(gpio.Low); err != nil {
		return err
	}
	if err := device.openPin.Out(gpio.High); err != nil {
		return err
	}
	device.state <- "opening"
	//TODO: handle error..
	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			// return err
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			// return err
		}
		device.state <- "stopped"
		// return nil
	})

	return nil
}

func (d *deviceController) CloseShutter(shutterID int64) error {
	device, err := d.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	if device.timer != nil {
		device.timer.Stop()
	}
	if err := device.openPin.Out(gpio.Low); err != nil {
		return err
	}
	if err := device.closePin.Out(gpio.High); err != nil {
		return err
	}
	device.state <- "closing"
	//TODO: handle error
	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			// return err
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			// return err
		}
		device.state <- "stopped"
		// return nil
	})
	return nil
}

func (d *deviceController) StopShutter(shutterID int64) error {
	device, err := d.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	if device.timer != nil {
		device.timer.Stop()
	}
	if err := device.openPin.Out(gpio.Low); err != nil {
		return err
	}
	if err := device.closePin.Out(gpio.Low); err != nil {
		return err
	}
	device.state <- "stopped"
	return nil
}

func (d *deviceController) TurnLightingOn(lightingID int64) error {
	device, err := d.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	if err := device.switchPin.Out(gpio.High); err != nil {
		return err
	}
	device.state <- "on"
	return nil
}

func (d *deviceController) TurnLightingOff(lightingID int64) error {
	device, err := d.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	if err := device.switchPin.Out(gpio.Low); err != nil {
		return err
	}
	device.state <- "off"
	return nil
}

func (d *deviceController) ScheduleShutterJobs(shutter *model.Shutter) error {
	device, err := d.getShutterByModel(shutter)
	if err != nil {
		return err
	}
	device.openJob, err = scheduler.Every().Day().NotImmediately().At(fmt.Sprintf("%02d:%02d", shutter.OpenTime.Hour(), shutter.OpenTime.Minute())).Run(func() {
		d.OpenShutter(device.ID)
	})
	if err != nil {
		return err
	}
	device.closeJob, err = scheduler.Every().Day().NotImmediately().At(fmt.Sprintf("%02d:%02d", shutter.CloseTime.Hour(), shutter.CloseTime.Minute())).Run(func() {
		d.CloseShutter(device.ID)
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *deviceController) UnscheduleShutterJobs(shutterID int64) error {
	device, err := d.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	if device.openJob != nil {
		device.openJob.Quit <- true
	}
	if device.closeJob != nil {
		device.closeJob.Quit <- true
	}
	return nil
}

func (d *deviceController) ScheduleLightingJobs(lighting *model.Lighting) error {
	device, err := d.getLightingByModel(lighting)
	if err != nil {
		return err
	}
	device.onJob, err = scheduler.Every().Day().NotImmediately().At(fmt.Sprintf("%02d:%02d", lighting.OnTime.Hour(), lighting.OnTime.Minute())).Run(func() {
		d.TurnLightingOn(device.ID)
	})
	if err != nil {
		return err
	}
	device.offJob, err = scheduler.Every().Day().NotImmediately().At(fmt.Sprintf("%02d:%02d", lighting.OffTime.Hour(), lighting.OffTime.Minute())).Run(func() {
		d.TurnLightingOff(device.ID)
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *deviceController) UnscheduleLightingJobs(lightingID int64) error {
	device, err := d.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	if device.onJob != nil {
		device.onJob.Quit <- true
	}
	if device.offJob != nil {
		device.offJob.Quit <- true
	}
	return nil
}

func (d *deviceController) getShutterByID(shutterID int64) (*shutter, error) {
	device, ok := d.shutters[shutterID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", shutterID)
	}
	return device, nil
}

func (d *deviceController) getShutterByModel(shutter *model.Shutter) (*shutter, error) {
	device, ok := d.shutters[shutter.ID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", shutter.ID)
	}
	return device, nil
}

func (d *deviceController) getLightingByID(lightingID int64) (*lighting, error) {
	device, ok := d.lightings[lightingID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", lightingID)
	}
	return device, nil
}

func (d *deviceController) getLightingByModel(lighting *model.Lighting) (*lighting, error) {
	device, ok := d.lightings[lighting.ID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", lighting.ID)
	}
	return device, nil
}
