package embedded

import (
	"fmt"
	"time"

	"github.com/he4d/almue/model"
	"github.com/he4d/scheduler"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type shutter struct {
	openPin             gpio.PinIO
	closePin            gpio.PinIO
	openJob             *scheduler.Job
	closeJob            *scheduler.Job
	completeWayDuration time.Duration
	timer               *time.Timer
}

func (c *EmbeddedController) RegisterShutters(shutters ...*model.Shutter) error {
	for _, shutterModel := range shutters {
		var openPin gpio.PinIO
		var closePin gpio.PinIO
		if c.simulate {
			openPin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.OpenPin}
			closePin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.ClosePin}
		} else {
			openPin = gpioreg.ByNumber(*shutterModel.OpenPin)
			closePin = gpioreg.ByNumber(*shutterModel.ClosePin)
		}
		duration := time.Duration(*shutterModel.CompleteWayInSeconds) * time.Second
		shutterToAdd := &shutter{
			openPin:             openPin,
			closePin:            closePin,
			completeWayDuration: duration,
		}

		c.shutters[shutterModel.ID] = shutterToAdd

		if shutterModel.JobsEnabled {
			if err := c.ScheduleShutterJobs(shutterModel); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EmbeddedController) UnregisterShutter(shutterID int64) error {
	if err := c.StopShutter(shutterID); err != nil {
		return err
	}

	if err := c.UnscheduleShutterJobs(shutterID); err != nil {
		return err
	}

	delete(c.shutters, shutterID)
	return nil
}

func (c *EmbeddedController) UpdateShutter(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
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
			c.UnregisterShutter(updatedShutter.ID)
			return nil
		}
		c.RegisterShutters(updatedShutter)
		return nil
	}
	if diffs.HasFlag(model.DIFFJOBSENABLED) {
		if updatedShutter.JobsEnabled {
			if err := c.ScheduleShutterJobs(updatedShutter); err != nil {
				return err
			}
			alreadyScheduled = true
		} else {
			if err := c.UnscheduleShutterJobs(updatedShutter.ID); err != nil {
				return err
			}
		}
	}
	if diffs.HasFlag(model.DIFFOPENPIN) || diffs.HasFlag(model.DIFFCLOSEPIN) {
		if err := c.changeShutterPins(diffs, updatedShutter); err != nil {
			return err
		}
	}
	if diffs.HasFlag(model.DIFFCOMPLETEWAYINSECONDS) {
		shutter, err := c.getShutterByID(updatedShutter.ID)
		if err != nil {
			return err
		}
		c.StopShutter(updatedShutter.ID)
		shutter.completeWayDuration = time.Duration(*updatedShutter.CompleteWayInSeconds) * time.Second
	}
	if diffs.HasFlag(model.DIFFOPENTIME) || diffs.HasFlag(model.DIFFCLOSETIME) {
		if updatedShutter.JobsEnabled && !alreadyScheduled {
			if err := c.rescheduleShutterJobs(updatedShutter); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EmbeddedController) OpenShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
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
	if err := c.stateStore.UpdateShutterState(shutterID, "opening"); err != nil {
		return err
	}
	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			//TODO: handle error..
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			//TODO: handle error..
		}
		if err := c.stateStore.UpdateShutterState(shutterID, "stopped"); err != nil {
			//TODO: handle error..
		}
	})

	return nil
}

func (c *EmbeddedController) CloseShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
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
	if err := c.stateStore.UpdateShutterState(shutterID, "closing"); err != nil {
		return err
	}

	device.timer = time.AfterFunc(device.completeWayDuration, func() {
		if err := device.closePin.Out(gpio.Low); err != nil {
			//TODO: handle error
		}
		if err := device.openPin.Out(gpio.Low); err != nil {
			//TODO: handle error
		}
		if err := c.stateStore.UpdateShutterState(shutterID, "stopped"); err != nil {
			//TODO: handle error..
		}
	})
	return nil
}

func (c *EmbeddedController) StopShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
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
	if err := c.stateStore.UpdateShutterState(shutterID, "stopped"); err != nil {
		return err
	}
	return nil
}

func (c *EmbeddedController) ScheduleShutterJobs(shutter *model.Shutter) error {
	device, err := c.getShutterByID(shutter.ID)
	if err != nil {
		return err
	}
	device.openJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", shutter.OpenTime.Hour(), shutter.OpenTime.Minute())).Run(func() {
		c.OpenShutter(shutter.ID)
	})
	if err != nil {
		return err
	}
	device.closeJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", shutter.CloseTime.Hour(), shutter.CloseTime.Minute())).Run(func() {
		c.CloseShutter(shutter.ID)
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *EmbeddedController) UnscheduleShutterJobs(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
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

func (c *EmbeddedController) getShutterByID(shutterID int64) (*shutter, error) {
	device, ok := c.shutters[shutterID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", shutterID)
	}
	return device, nil
}

func (c *EmbeddedController) rescheduleShutterJobs(shutter *model.Shutter) error {
	if err := c.UnscheduleShutterJobs(shutter.ID); err != nil {
		return err
	}
	if err := c.ScheduleShutterJobs(shutter); err != nil {
		return err
	}
	return nil
}

func (c *EmbeddedController) changeShutterPins(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
	shutter, err := c.getShutterByID(updatedShutter.ID)
	if err != nil {
		return err
	}
	c.StopShutter(updatedShutter.ID)
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if c.simulate {
			shutter.openPin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.OpenPin}
		} else {
			shutter.openPin = gpioreg.ByNumber(*updatedShutter.OpenPin)
		}
	}
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if err != nil {
			return err
		}
		if c.simulate {
			shutter.closePin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.ClosePin}
		} else {
			shutter.closePin = gpioreg.ByNumber(*updatedShutter.ClosePin)
		}
	}
	return nil
}
