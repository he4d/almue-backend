package model

import (
	"time"
)

//Shutter represents the database object of a shutter
type Shutter struct {
	Base
	Description          *string   `json:"description"`
	OpenPin              *int      `json:"openPin"`
	ClosePin             *int      `json:"closePin"`
	CompleteWayInSeconds *int      `json:"completeWayInSeconds"`
	JobsEnabled          bool      `json:"jobsEnabled"`
	OpenTime             time.Time `json:"openTime"`
	CloseTime            time.Time `json:"closeTime"`
	EmergencyEnabled     bool      `json:"emergencyEnabled"`
	DeviceStatus         string    `json:"deviceStatus"`
	Disabled             bool      `json:"disabled"`
	FloorID              int64     `json:"floorId"`
}
