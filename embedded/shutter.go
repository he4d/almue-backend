package embedded

import (
	"fmt"
	"strconv"
	"time"

	"sync"

	"github.com/he4d/almue-backend/model"
	"github.com/he4d/scheduler"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type shutter struct {
	sync.Mutex
	openPin             gpio.PinIO
	closePin            gpio.PinIO
	openJob             *scheduler.Job
	closeJob            *scheduler.Job
	completeWayDuration time.Duration
	timer               *time.Timer
	ticker              *time.Ticker
	openingInPrc        int
}

func (s *shutter) getTickDuration() time.Duration {
	calc := (s.completeWayDuration.Seconds() * 5.0 / 100.0) * 1000.0
	return time.Millisecond * time.Duration(calc)
}

// RegisterShutters registers one or more shutters to the controller. It will also start the scheudle if enabled for the given shutter
func (c *Controller) RegisterShutters(shutters ...*model.Shutter) error {
	for _, shutterModel := range shutters {
		var openPin gpio.PinIO
		var closePin gpio.PinIO
		if c.simulate {
			openPin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.OpenPin}
			closePin = &simulatePinIO{name: *shutterModel.Description, number: *shutterModel.ClosePin}
		} else {
			openPin = gpioreg.ByName(strconv.Itoa(*shutterModel.OpenPin))
			closePin = gpioreg.ByName(strconv.Itoa(*shutterModel.ClosePin))
		}
		duration := time.Duration(*shutterModel.CompleteWayInSeconds) * time.Second
		shutterToAdd := &shutter{
			openPin:             openPin,
			closePin:            closePin,
			completeWayDuration: duration,
			openingInPrc:        shutterModel.OpeningInPrc,
		}

		c.shuttersLock.Lock()
		c.shutters[shutterModel.ID] = shutterToAdd
		c.shuttersLock.Unlock()

		if shutterModel.JobsEnabled {
			if err := c.ScheduleShutterJobs(shutterModel); err != nil {
				return err
			}
		}
	}
	return nil
}

// UnregisterShutter unregisters the shutter with the given id from the controller
func (c *Controller) UnregisterShutter(shutterID int64) error {
	if err := c.StopShutter(shutterID); err != nil {
		return err
	}

	if err := c.UnscheduleShutterJobs(shutterID); err != nil {
		return err
	}

	c.shuttersLock.Lock()
	delete(c.shutters, shutterID)
	c.shuttersLock.Unlock()
	return nil
}

// UpdateShutter updates a Shutter according to the differences that get passed
func (c *Controller) UpdateShutter(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
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

// OpenShutter opens the shutter with the given id
// It also updates the state store
func (c *Controller) OpenShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	device.Lock()
	defer device.Unlock()
	if device.ticker != nil {
		device.ticker.Stop()
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
	if device.openingInPrc == 100.0 {
		// REFERENCE DRIVE
		if err := c.stateStore.UpdateShutterState(shutterID, "referencing"); err != nil {
			return err
		}
		device.timer = time.AfterFunc(device.completeWayDuration, func() {
			if err := c.StopShutter(shutterID); err != nil {
				//TODO: Handle error
			}
		})
	} else {
		// NORMAL DRIVE
		if err := c.stateStore.UpdateShutterState(shutterID, "opening"); err != nil {
			return err
		}
		device.ticker = time.NewTicker(device.getTickDuration())
		go func() {
			for range device.ticker.C {
				device.openingInPrc += 5
				if err := c.stateStore.UpdateShutterOpening(shutterID, device.openingInPrc); err != nil {
					//TODO: Handle error
				}
				if device.openingInPrc == 100 {
					if err := c.StopShutter(shutterID); err != nil {
						//TODO: Handle error
					}
					device.ticker.Stop()
					device.ticker = nil
					return
				}
			}
		}()
	}
	return nil
}

// CloseShutter closes the shutter with the given id
// It also updates the state store
func (c *Controller) CloseShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	device.Lock()
	defer device.Unlock()
	if device.ticker != nil {
		device.ticker.Stop()
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
	if device.openingInPrc == 0 {
		if err := c.stateStore.UpdateShutterState(shutterID, "referencing"); err != nil {
			return err
		}
		device.timer = time.AfterFunc(device.completeWayDuration, func() {
			if err := c.StopShutter(shutterID); err != nil {
				//TODO: Handle error
			}
		})
	} else {
		// NORMAL DRIVE
		if err := c.stateStore.UpdateShutterState(shutterID, "closing"); err != nil {
			return err
		}
		device.ticker = time.NewTicker(device.getTickDuration())
		go func() {
			for range device.ticker.C {
				device.openingInPrc -= 5
				if err := c.stateStore.UpdateShutterOpening(shutterID, device.openingInPrc); err != nil {
					//TODO: Handle error
				}
				if device.openingInPrc == 0 {
					if err := c.StopShutter(shutterID); err != nil {
						//TODO: Handle error
					}
					device.ticker.Stop()
					device.ticker = nil
					return
				}
			}
		}()
	}
	return nil
}

// StopShutter stops the shutter with the given id
// It also updates the state store
func (c *Controller) StopShutter(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	device.Lock()
	defer device.Unlock()
	if device.ticker != nil {
		device.ticker.Stop()
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

// ScheduleShutterJobs schedules jobs of the given shutter
func (c *Controller) ScheduleShutterJobs(shutter *model.Shutter) error {
	device, err := c.getShutterByID(shutter.ID)
	if err != nil {
		return err
	}
	device.Lock()
	defer device.Unlock()
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

// UnscheduleShutterJobs unschedules jobs of the given shutter
func (c *Controller) UnscheduleShutterJobs(shutterID int64) error {
	device, err := c.getShutterByID(shutterID)
	if err != nil {
		return err
	}
	device.Lock()
	if device.openJob != nil {
		device.openJob.Quit <- true
	}
	if device.closeJob != nil {
		device.closeJob.Quit <- true
	}
	device.Unlock()
	return nil
}

func (c *Controller) getShutterByID(shutterID int64) (*shutter, error) {
	c.shuttersLock.RLock()
	device, ok := c.shutters[shutterID]
	c.shuttersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", shutterID)
	}
	return device, nil
}

func (c *Controller) rescheduleShutterJobs(shutter *model.Shutter) error {
	if err := c.UnscheduleShutterJobs(shutter.ID); err != nil {
		return err
	}
	if err := c.ScheduleShutterJobs(shutter); err != nil {
		return err
	}
	return nil
}

func (c *Controller) changeShutterPins(diffs model.DifferenceType, updatedShutter *model.Shutter) error {
	c.StopShutter(updatedShutter.ID)
	shutter, err := c.getShutterByID(updatedShutter.ID)
	if err != nil {
		return err
	}
	shutter.Lock()
	defer shutter.Unlock()
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if c.simulate {
			shutter.openPin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.OpenPin}
		} else {
			shutter.openPin = gpioreg.ByName(strconv.Itoa(*updatedShutter.OpenPin))
		}
	}
	if diffs.HasFlag(model.DIFFOPENPIN) {
		if err != nil {
			return err
		}
		if c.simulate {
			shutter.closePin = &simulatePinIO{name: *updatedShutter.Description, number: *updatedShutter.ClosePin}
		} else {
			shutter.closePin = gpioreg.ByName(strconv.Itoa(*updatedShutter.ClosePin))
		}
	}
	return nil
}
