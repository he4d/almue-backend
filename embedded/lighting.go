package embedded

import "periph.io/x/periph/conn/gpio"
import "github.com/carlescere/scheduler"

type lighting struct {
	ID        int64
	switchPin gpio.PinIO
	onJob     *scheduler.Job
	offJob    *scheduler.Job
	state     chan string
	quit      chan bool
}