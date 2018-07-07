package embedded

import (
	"fmt"
	"log"

	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
)

type simulatePinIO struct {
	name   string
	number int
}

func (s *simulatePinIO) String() string {
	return fmt.Sprintf("%s(%d)", s.name, s.number)
}

func (s *simulatePinIO) Name() string {
	return s.name
}

func (s *simulatePinIO) Number() int {
	return s.number
}

func (s *simulatePinIO) Function() string {
	return "Out"
}

func (s *simulatePinIO) In(pull gpio.Pull, edge gpio.Edge) error {
	//TODO:
	return nil
}

func (s *simulatePinIO) Read() gpio.Level {
	//TODO:
	return true
}

func (s *simulatePinIO) WaitForEdge(timeout time.Duration) bool {
	//TODO:
	return true
}

func (s *simulatePinIO) Pull() gpio.Pull {
	//TODO:
	return 2
}

func (s *simulatePinIO) DefaultPull() gpio.Pull {
	//TODO:
	return 2
}

func (s *simulatePinIO) Halt() error {
	//TODO:
	return nil
}

func (s *simulatePinIO) PWM(gpio.Duty, physic.Frequency) error {
	//TODO:
	return nil
}

func (s *simulatePinIO) Out(l gpio.Level) error {
	log.Printf("Name: %s - Pin: %d is switching level to %t\n", s.name, s.number, l)
	return nil
}
