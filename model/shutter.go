package model

import (
	"time"
)

//Shutter represents the database object of a shutter
type Shutter struct {
	Base                 `valid:"-"`
	Description          string    `json:"description" valid:"alphanum,required"`
	OpenPin              int       `json:"openPin" valid:"gpio,required"`
	ClosePin             int       `json:"closePin" valid:"gpio,required"`
	CompleteWayInSeconds int       `json:"completeWayInSeconds" valid:"required"`
	TimerEnabled         bool      `json:"timerEnabled" valid:"-"`
	OpenTime             time.Time `json:"openTime" valid:"-"`
	CloseTime            time.Time `json:"closeTime" valid:"-"`
	EmergencyEnabled     bool      `json:"emergencyEnabled" valid:"-"`
	DeviceStatus         string    `json:"deviceStatus" valid:"-"`
	Disabled             bool      `json:"disabled" valid:"-"`
	FloorID              int64     `json:"floorId" valid:"-"`
}
