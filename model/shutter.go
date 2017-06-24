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

//DeepCopy creates a deep copy of a Shutter
func (s *Shutter) DeepCopy() *Shutter {
	if s == nil {
		return nil
	}
	descr := *s.Description
	openPin := *s.OpenPin
	closePin := *s.ClosePin
	completeWayInSecs := *s.CompleteWayInSeconds
	copy := &Shutter{
		Base:                 s.Base,
		Description:          &descr,
		OpenPin:              &openPin,
		ClosePin:             &closePin,
		CompleteWayInSeconds: &completeWayInSecs,
		JobsEnabled:          s.JobsEnabled,
		OpenTime:             s.OpenTime,
		CloseTime:            s.CloseTime,
		EmergencyEnabled:     s.EmergencyEnabled,
		DeviceStatus:         s.DeviceStatus,
		Disabled:             s.Disabled,
		FloorID:              s.FloorID,
	}
	return copy
}
