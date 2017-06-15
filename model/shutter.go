package model

import (
	"errors"
	"time"
)

type Shutter struct {
	ModelBase
	Description          string    `json:"description"`
	OpenPin              int       `json:"openPin"`
	ClosePin             int       `json:"closePin"`
	CompleteWayInSeconds int       `json:"completeWayInSeconds"`
	TimerEnabled         bool      `json:"timerEnabled"`
	OpenTime             time.Time `json:"openTime"`
	CloseTime            time.Time `json:"closeTime"`
	EmergencyEnabled     bool      `json:"emergencyEnabled"`
	DeviceStatus         int       `json:"deviceStatus"`
	Disabled             bool      `json:"disabled"`
	FloorID              int64     `json:"floorId"`
}

func (s Shutter) Validate() error {
	return errors.New("not implemented")
}
