package embedded

import (
	"errors"
	"log"
	"time"

	"periph.io/x/periph/host"

	"fmt"

	"github.com/carlescere/scheduler"
	"github.com/he4d/almue/model"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
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

func (d *deviceController) GetShutterStateSyncChannels(shutterID int64) (*StateSyncChannels, error) {
	shutter, ok := d.shutters[shutterID]
	if !ok {
		return nil, errors.New("Could not obtain shutter for getting the statesyncchannels")
	}
	return shutter.stateSync, nil
}

func (d *deviceController) GetLightingStateSyncChannels(lightingID int64) (*StateSyncChannels, error) {
	lighting, ok := d.lightings[lightingID]
	if !ok {
		return nil, errors.New("Could not obtain lighting for getting the statesyncchannels")
	}
	return lighting.stateSync, nil
}

func (d *deviceController) RegisterShutters(shutters ...*model.Shutter) error {
	for _, shutterModel := range shutters {
		var openPin gpio.PinIO
		var closePin gpio.PinIO
		if d.simulate {
			openPin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.OpenPin}
			closePin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.ClosePin}
		} else {
			openPin = gpioreg.ByNumber(*shutterModel.OpenPin)
			closePin = gpioreg.ByNumber(*shutterModel.ClosePin)
		}
		duration := time.Duration(*shutterModel.CompleteWayInSeconds) * time.Second
		stateSync := new(StateSyncChannels)
		stateSync.State = make(chan string)
		stateSync.Quit = make(chan bool)
		d.shutters[shutterModel.ID] = &shutter{
			ID:                  shutterModel.ID,
			openPin:             openPin,
			closePin:            closePin,
			completeWayDuration: duration,
			stateSync:           stateSync,
		}
		if shutterModel.JobsEnabled {
			if err := d.ScheduleShutterJobs(shutterModel); err != nil {
				return err
			}
		}
	}
	return nil
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
	shutter.stateSync.Quit <- true

	delete(d.shutters, shutterID)
	return nil
}

func (d *deviceController) UpdateShutter(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
	if diffs.HasFlag(model.DIFFNONE) {
		return nil
	}
	if diffs.HasFlag(model.DIFFEMERGENCYENABLED) {
		//TODO: Emergencydevices...
		return nil
	}
	if diffs.HasFlag(model.DIFFDISABLED) {
		if updatedShutter.Disabled {
			d.UnregisterShutter(updatedShutter.ID)
			return nil
		}
		d.RegisterShutters(updatedShutter)
		return nil
	}
	if diffs.HasFlag(model.DIFFJOBSENABLED) {
		if updatedShutter.JobsEnabled {
			if err := d.ScheduleShutterJobs(updatedShutter); err != nil {
				return err
			}
		} else {
			if err := d.UnscheduleShutterJobs(updatedShutter.ID); err != nil {
				return err
			}
		}
	}
	if diffs.HasFlag(model.DIFFOPENPIN) || diffs.HasFlag(model.DIFFCLOSEPIN) {
		if err := d.changeShutterPins(diffs, updatedShutter); err != nil {
			return err
		}
	}
	if diffs.HasFlag(model.DIFFCOMPLETEWAYINSECONDS) {
		shutter, err := d.getShutterByID(updatedShutter.ID)
		if err != nil {
			return err
		}
		shutter.completeWayDuration = time.Duration(*updatedShutter.CompleteWayInSeconds) * time.Second
	}
	if diffs.HasFlag(model.DIFFOPENTIME) || diffs.HasFlag(model.DIFFCLOSETIME) {
		if err := d.rescheduleShutterJobs(updatedShutter); err != nil {
			return err
		}
	}
	return nil
}

func (d *deviceController) RegisterLightings(lightings ...*model.Lighting) error {
	for _, lightingModel := range lightings {
		var switchPin gpio.PinIO
		if d.simulate {
			switchPin = &simulatePinIO{name: *lightingModel.Description, number: *lightingModel.SwitchPin}
		} else {
			switchPin = gpioreg.ByNumber(*lightingModel.SwitchPin)
		}
		stateSync := new(StateSyncChannels)
		stateSync.State = make(chan string)
		stateSync.Quit = make(chan bool)
		d.lightings[lightingModel.ID] = &lighting{
			ID:        lightingModel.ID,
			switchPin: switchPin,
			stateSync: stateSync,
		}
		if lightingModel.JobsEnabled {
			if err := d.ScheduleLightingJobs(lightingModel); err != nil {
				return err
			}
		}
	}
	return nil
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
	lighting.stateSync.Quit <- true

	delete(d.lightings, lightingID)
	return nil
}

func (d *deviceController) UpdateLighting(diffs model.DifferenceType, updatedLighting *model.Lighting) error {
	if diffs.HasFlag(model.DIFFNONE) {
		return nil
	}
	if diffs.HasFlag(model.DIFFEMERGENCYENABLED) {
		//TODO: Emergencydevices...
		return nil
	}
	if diffs.HasFlag(model.DIFFDISABLED) {
		if updatedLighting.Disabled {
			d.UnregisterLighting(updatedLighting.ID)
			return nil
		}
		d.RegisterLightings(updatedLighting)
		return nil
	}
	if diffs.HasFlag(model.DIFFJOBSENABLED) {
		if updatedLighting.JobsEnabled {
			if err := d.ScheduleLightingJobs(updatedLighting); err != nil {
				return err
			}
		} else {
			if err := d.UnscheduleLightingJobs(updatedLighting.ID); err != nil {
				return err
			}
		}
	}
	if diffs.HasFlag(model.DIFFSWITCHPIN) {
		if err := d.changeLightingPin(diffs, updatedLighting); err != nil {
			return err
		}
	}
	if diffs.HasFlag(model.DIFFONTIME) || diffs.HasFlag(model.DIFFOFFTIME) {
		if err := d.rescheduleLightingJobs(updatedLighting); err != nil {
			return err
		}
	}
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
	device.stateSync.State <- "opening"
	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			//TODO: handle error..
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			//TODO: handle error..
		}
		device.stateSync.State <- "stopped"
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
	device.stateSync.State <- "closing"

	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			//TODO: handle error
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			//TODO: handle error
		}
		device.stateSync.State <- "stopped"
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
	device.stateSync.State <- "stopped"
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
	device.stateSync.State <- "on"
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
	device.stateSync.State <- "off"
	return nil
}

func (d *deviceController) ScheduleShutterJobs(shutter *model.Shutter) error {
	device, err := d.getShutterByID(shutter.ID)
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
	device, err := d.getLightingByID(lighting.ID)
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

func (d *deviceController) changeShutterPins(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
	return nil
}

func (d *deviceController) changeLightingPin(diffs model.DifferenceType, updatedLighting *model.Lighting) error {
	return nil
}

func (d *deviceController) rescheduleShutterJobs(shutter *model.Shutter) error {
	if err := d.UnscheduleShutterJobs(shutter.ID); err != nil {
		return err
	}
	if err := d.ScheduleShutterJobs(shutter); err != nil {
		return err
	}
	return nil
}

func (d *deviceController) rescheduleLightingJobs(lighting *model.Lighting) error {
	if err := d.UnscheduleLightingJobs(lighting.ID); err != nil {
		return err
	}
	if err := d.ScheduleLightingJobs(lighting); err != nil {
		return err
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

func (d *deviceController) getLightingByID(lightingID int64) (*lighting, error) {
	device, ok := d.lightings[lightingID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", lightingID)
	}
	return device, nil
}
