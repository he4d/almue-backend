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

//DeepCopy creates a deep copy of a Lighting
func (l *Lighting) DeepCopy() *Lighting {
	if l == nil {
		return nil
	}
	descr := *l.Description
	switchPin := *l.SwitchPin
	copy := &Lighting{
		Base:             l.Base,
		Description:      &descr,
		SwitchPin:        &switchPin,
		JobsEnabled:      l.JobsEnabled,
		OnTime:           l.OnTime,
		OffTime:          l.OffTime,
		EmergencyEnabled: l.EmergencyEnabled,
		DeviceStatus:     l.DeviceStatus,
		Disabled:         l.Disabled,
		FloorID:          l.FloorID,
	}
	return copy
}
