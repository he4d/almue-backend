package embedded

import (
	"errors"
	"fmt"
	"time"

	"github.com/he4d/almue/model"
	"github.com/he4d/scheduler"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type shutter struct {
	ID                  int64
	openPin             gpio.PinIO
	closePin            gpio.PinIO
	openJob             *scheduler.Job
	closeJob            *scheduler.Job
	completeWayDuration time.Duration
	timer               *time.Timer
	stateSync           *StateSyncChannels
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

	if err := d.StopShutter(shutterID); err != nil {
		return err
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
	if diffs == model.DIFFNONE {
		return nil
	}

	var alreadyScheduled bool

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
			alreadyScheduled = true
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
		d.StopShutter(updatedShutter.ID)
		shutter.completeWayDuration = time.Duration(*updatedShutter.CompleteWayInSeconds) * time.Second
	}
	if diffs.HasFlag(model.DIFFOPENTIME) || diffs.HasFlag(model.DIFFCLOSETIME) {
		if updatedShutter.JobsEnabled && !alreadyScheduled {
			if err := d.rescheduleShutterJobs(updatedShutter); err != nil {
				return err
			}
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

func (d *deviceController) ScheduleShutterJobs(shutter *model.Shutter) error {
	device, err := d.getShutterByID(shutter.ID)
	if err != nil {
		return err
	}
	device.openJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", shutter.OpenTime.Hour(), shutter.OpenTime.Minute())).Run(func() {
		d.OpenShutter(device.ID)
	})
	if err != nil {
		return err
	}
	device.closeJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", shutter.CloseTime.Hour(), shutter.CloseTime.Minute())).Run(func() {
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

func (d *deviceController) GetShutterStateSyncChannels(shutterID int64) (*StateSyncChannels, error) {
	shutter, ok := d.shutters[shutterID]
	if !ok {
		return nil, errors.New("Could not obtain shutter for getting the statesyncchannels")
	}
	return shutter.stateSync, nil
}

func (d *deviceController) getShutterByID(shutterID int64) (*shutter, error) {
	device, ok := d.shutters[shutterID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", shutterID)
	}
	return device, nil
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

func (d *deviceController) changeShutterPins(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
	shutter, err := d.getShutterByID(updatedShutter.ID)
	if err != nil {
		return err
	}
	d.StopShutter(updatedShutter.ID)
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if d.simulate {
			shutter.openPin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.OpenPin}
		} else {
			shutter.openPin = gpioreg.ByNumber(*updatedShutter.OpenPin)
		}
	}
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if err != nil {
			return err
		}
		if d.simulate {
			shutter.closePin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.ClosePin}
		} else {
			shutter.closePin = gpioreg.ByNumber(*updatedShutter.ClosePin)
		}
	}
	return nil
}
