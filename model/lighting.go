package model

import (
	"time"
)

//Lighting represents the database object of a lighting
type Lighting struct {
	Base
	Description      *string   `json:"description"`
	SwitchPin        *int      `json:"switchPin"`
	JobsEnabled      bool      `json:"jobsEnabled"`
	OnTime           time.Time `json:"onTime"`
	OffTime          time.Time `json:"offTime"`
	EmergencyEnabled bool      `json:"emergencyEnabled"`
	DeviceStatus     string    `json:"deviceStatus"`
	Disabled         bool      `json:"disabled"`
	FloorID          int64     `json:"floorId"`
}
