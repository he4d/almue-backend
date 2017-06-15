package model

import (
	"errors"
	"time"
)

type Lighting struct {
	ModelBase
	Description      string    `json:"description"`
	SwitchPin        int       `json:"switchPin"`
	TimerEnabled     bool      `json:"timerEnabled"`
	OnTime           time.Time `json:"onTime"`
	OffTime          time.Time `json:"offTime"`
	EmergencyEnabled bool      `json:"emergencyEnabled"`
	DeviceStatus     int       `json:"deviceStatus"`
	Disabled         bool      `json:"disabled"`
	FloorID          int64     `json:"floorId"`
}

func (l Lighting) Validate() error {
	return errors.New("not implemented")
}
