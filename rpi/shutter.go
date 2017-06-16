package rpi

import (
	"github.com/carlescere/scheduler"
	"periph.io/x/periph/conn/gpio"
)

type shutter struct {
	ID       int64
	openPin  gpio.PinIO
	closePin gpio.PinIO
	openJob  *scheduler.Job
	closeJob *scheduler.Job
}
