package model

import (
	"time"
)

//Lighting represents the database object of a lighting
type Lighting struct {
	Base             `valid:"-"`
	Description      string    `json:"description" valid:"alphanum,required"`
	SwitchPin        int       `json:"switchPin" valid:"gpio,required"`
	TimerEnabled     bool      `json:"timerEnabled" valid:"-"`
	OnTime           time.Time `json:"onTime" valid:"-"`
	OffTime          time.Time `json:"offTime" valid:"-"`
	EmergencyEnabled bool      `json:"emergencyEnabled" valid:"-"`
	DeviceStatus     int       `json:"deviceStatus" valid:"int"`
	Disabled         bool      `json:"disabled" valid:"-"`
	FloorID          int64     `json:"floorId" valid:"-"`
}
