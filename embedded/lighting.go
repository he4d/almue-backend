package embedded

import (
	"errors"
	"fmt"

	"github.com/he4d/almue/model"
	"github.com/he4d/scheduler"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type lighting struct {
	ID        int64
	switchPin gpio.PinIO
	onJob     *scheduler.Job
	offJob    *scheduler.Job
	stateSync *StateSyncChannels
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
	if err := d.TurnLightingOff(lightingID); err != nil {
		return err
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
			alreadyScheduled = true
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
		if updatedLighting.JobsEnabled && !alreadyScheduled {
			if err := d.rescheduleLightingJobs(updatedLighting); err != nil {
				return err
			}
		}
	}
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

func (d *deviceController) ScheduleLightingJobs(lighting *model.Lighting) error {
	device, err := d.getLightingByID(lighting.ID)
	if err != nil {
		return err
	}
	device.onJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", lighting.OnTime.Hour(), lighting.OnTime.Minute())).Run(func() {
		d.TurnLightingOn(device.ID)
	})
	if err != nil {
		return err
	}
	device.offJob, err = scheduler.Every().Day().At(fmt.Sprintf("%02d:%02d", lighting.OffTime.Hour(), lighting.OffTime.Minute())).Run(func() {
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

func (d *deviceController) GetLightingStateSyncChannels(lightingID int64) (*StateSyncChannels, error) {
	lighting, ok := d.lightings[lightingID]
	if !ok {
		return nil, errors.New("Could not obtain lighting for getting the statesyncchannels")
	}
	return lighting.stateSync, nil
}

func (d *deviceController) changeLightingPin(diffs model.DifferenceType, updatedLighting *model.Lighting) error {
	d.TurnLightingOff(updatedLighting.ID)
	lighting, err := d.getLightingByID(updatedLighting.ID)
	if err != nil {
		return err
	}
	if d.simulate {
		lighting.switchPin = &simulatePinIO{name: *updatedLighting.Description, number: *updatedLighting.SwitchPin}
	} else {
		lighting.switchPin = gpioreg.ByNumber(*updatedLighting.SwitchPin)
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

func (d *deviceController) getLightingByID(lightingID int64) (*lighting, error) {
	device, ok := d.lightings[lightingID]
	if !ok {
		return nil, fmt.Errorf("Device with ID: %d is not registered in the DeviceController", lightingID)
	}
	return device, nil
}
