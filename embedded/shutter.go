package embedded

import (
	"time"

	"github.com/carlescere/scheduler"
	"periph.io/x/periph/conn/gpio"
)

type shutter struct {
	ID                  int64
	openPin             gpio.PinIO
	closePin            gpio.PinIO
	openJob             *scheduler.Job
	closeJob            *scheduler.Job
	completeWayDuration time.Duration
	timer               *time.Timer
	state               chan string
	quit                chan bool
}
