package embedded

import (
	"fmt"

	"sync"

	"github.com/he4d/almue/model"
	"github.com/he4d/scheduler"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type lighting struct {
	sync.Mutex
	switchPin gpio.PinIO
	onJob     *scheduler.Job
	offJob    *scheduler.Job
}

func (c *EmbeddedController) RegisterLightings(lightings ...*model.Lighting) error {
	for _, lightingModel := range lightings {
		var switchPin gpio.PinIO
		if c.simulate {
			switchPin = &simulatePinIO{name: *lightingModel.Description, number: *lightingModel.SwitchPin}
		} else {
			switchPin = gpioreg.ByNumber(*lightingModel.SwitchPin)
		}
		lightingToAdd := &lighting{
			switchPin: switchPin,
		}

		c.lightingsLock.Lock()
		c.lightings[lightingModel.ID] = lightingToAdd
		c.lightingsLock.Unlock()

		if lightingModel.JobsEnabled {
			if err := c.ScheduleLightingJobs(lightingModel); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EmbeddedController) UnregisterLighting(lightingID int64) error {
	if err := c.TurnLightingOff(lightingID); err != nil {
		return err
	}
	if err := c.UnscheduleLightingJobs(lightingID); err != nil {
		return err
	}

	c.lightingsLock.Lock()
	delete(c.lightings, lightingID)
	c.lightingsLock.Unlock()
	return nil
}

func (c *EmbeddedController) UpdateLighting(diffs model.DifferenceType, updatedLighting *model.Lighting) error {
	if diffs == model.DIFFNONE {
		return nil
	}

	var alreadyScheduled bool

	if diffs.HasFlag(model.DIFFEMERGENCYENABLED) {
		//TODO: Emergencydevices...
		return nil
	}
	if diffs.HasFlag(model.DIFFDISABLED) {
		if updatedLighting.Disabled {
			c.UnregisterLighting(updatedLighting.ID)
			return nil
		}
		c.RegisterLightings(updatedLighting)
		return nil
	}
	if diffs.HasFlag(model.DIFFJOBSENABLED) {
		if updatedLighting.JobsEnabled {
			if err := c.ScheduleLightingJobs(updatedLighting); err != nil {
				return err
			}
			alreadyScheduled = true
		} else {
			if err := c.UnscheduleLightingJobs(updatedLighting.ID); err != nil {
				return err
			}
		}
	}
	if diffs.HasFlag(model.DIFFSWITCHPIN) {
		if err := c.changeLightingPin(diffs, updatedLighting); err != nil {
			return err
		}
	}
	if diffs.HasFlag(model.DIFFONTIME) || diffs.HasFlag(model.DIFFOFFTIME) {
		if updatedLighting.JobsEnabled && !alreadyScheduled {
			if err := c.rescheduleLightingJobs(updatedLighting); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EmbeddedController) TurnLightingOn(lightingID int64) error {
	device, err := c.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	device.Lock()
	if err := device.switchPin.Out(gpio.High); err != nil {
		return err
	}
	if err := c.stateStore.UpdateLightingState(lightingID, "on"); err != nil {
		return err
	}
	device.Unlock()
	return nil
}

func (c *EmbeddedController) TurnLightingOff(lightingID int64) error {
	device, err := c.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	device.Lock()
	if err := device.switchPin.Out(gpio.Low); err != nil {
		return err
	}
	if err := c.stateStore.UpdateLightingState(lightingID, "off"); err != nil {
		return err
	}
	device.Unlock()
	return nil
}

func (c *EmbeddedController) ScheduleLightingJobs(lighting *model.Lighting) error {
	device, err := c.getLightingByID(lighting.ID)
	if err != nil {
		return err
	}
	device.Lock()
	device.onJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", lighting.OnTime.Hour(), lighting.OnTime.Minute())).Run(func() {
		c.TurnLightingOn(lighting.ID)
	})
	if err != nil {
		return err
	}
	device.offJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", lighting.OffTime.Hour(), lighting.OffTime.Minute())).Run(func() {
		c.TurnLightingOff(lighting.ID)
	})
	if err != nil {
		return err
	}
	device.Unlock()
	return nil
}

func (c *EmbeddedController) UnscheduleLightingJobs(lightingID int64) error {
	device, err := c.getLightingByID(lightingID)
	if err != nil {
		return err
	}
	device.Lock()
	if device.onJob != nil {
		device.onJob.Quit <- true
	}
	if device.offJob != nil {
		device.offJob.Quit <- true
	}
	device.Unlock()
	return nil
}

func (c *EmbeddedController) changeLightingPin(diffs model.DifferenceType, updatedLighting *model.Lighting) error {
	c.TurnLightingOff(updatedLighting.ID)
	lighting, err := c.getLightingByID(updatedLighting.ID)
	if err != nil {
		return err
	}
	lighting.Lock()
	if c.simulate {
		lighting.switchPin = &simulatePinIO{name: *updatedLighting.Description, number: *updatedLighting.SwitchPin}
	} else {
		lighting.switchPin = gpioreg.ByNumber(*updatedLighting.SwitchPin)
	}
	lighting.Unlock()
	return nil
}

func (c *EmbeddedController) rescheduleLightingJobs(lighting *model.Lighting) error {
	if err := c.UnscheduleLightingJobs(lighting.ID); err != nil {
		return err
	}
	if err := c.ScheduleLightingJobs(lighting); err != nil {
		return err
	}
	return nil
}

func (c *EmbeddedController) getLightingByID(lightingID int64) (*lighting, error) {
	c.lightingsLock.RLock()
	device, ok := c.lightings[lightingID]
	c.lightingsLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", lightingID)
	}
	return device, nil
}
