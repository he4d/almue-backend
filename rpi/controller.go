package rpi

import (
	"errors"
	"log"

	"github.com/he4d/almue/model"
	"github.com/kidoman/embd"
)

//DeviceController must be implemented by every controller that should control the devices
type DeviceController interface {
	RegisterShutters(shutters ...*model.Shutter) error

	RegisterLightings(lightings ...*model.Lighting) error

	OpenShutter(shutter *model.Shutter) error

	CloseShutter(shutter *model.Shutter) error

	StopShutter(shutter *model.Shutter) error

	TurnLightingOn(lighting *model.Lighting) error

	TurnLightingOff(lighting *model.Lighting) error
}

type deviceController struct {
	shutters  map[int64]*shutter
	lightings map[int64]*lighting
}

//New creates a new DeviceController and returns it
func New() DeviceController {
	err := embd.InitGPIO()
	if err != nil {
		log.Println(err)
		log.Fatalln("Getting new DeviceController failed")
	}
	return &deviceController{}
}

func (d *deviceController) RegisterShutters(shutters ...*model.Shutter) error {
	for _, shutterModel := range shutters {
		shutterID := shutterModel.ID
		openPin, err := embd.NewDigitalPin(shutterModel.OpenPin)
		if err != nil {
			return err
		}
		closePin, err := embd.NewDigitalPin(shutterModel.ClosePin)
		if err != nil {
			return err
		}
		d.shutters[shutterID] = &shutter{
			deviceID: shutterID,
			openPin:  &openPin,
			closePin: &closePin,
		}
	}
	return nil
}

func (d *deviceController) RegisterLightings(lightings ...*model.Lighting) error {
	return errors.New("not implemented")
}

func (d *deviceController) OpenShutter(shutter *model.Shutter) error {
	return errors.New("not implemented")
}

func (d *deviceController) CloseShutter(shutter *model.Shutter) error {
	return errors.New("not implemented")
}

func (d *deviceController) StopShutter(shutter *model.Shutter) error {
	return errors.New("not implemented")
}

func (d *deviceController) TurnLightingOn(lighting *model.Lighting) error {
	return errors.New("not implemented")
}

func (d *deviceController) TurnLightingOff(lighting *model.Lighting) error {
	return errors.New("not implemented")
}
