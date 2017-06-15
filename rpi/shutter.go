package rpi

import "github.com/kidoman/embd"

type shutter struct {
	deviceID int64
	openPin  *embd.DigitalPin
	closePin *embd.DigitalPin
}
