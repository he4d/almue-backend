package rpi

import (
	"github.com/kidoman/embd"
)

type lighting struct {
	deviceID  int64
	switchPin *embd.DigitalPin
}
