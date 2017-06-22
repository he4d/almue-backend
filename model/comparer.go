package model

type ModelDifference uint16

const (
	NONE ModelDifference = 1 << iota
	//DEVICE SPECIFIC
	EMERGENCYENABLED
	DISABLED
	JOBSENABLED
	//SHUTTER SPECIFIC
	OPENPIN
	CLOSEPIN
	COMPLETEWAYINSECONDS
	OPENTIME
	CLOSETIME
	//LIGHTING SPECIFIC
	SWITCHPIN
	ONTIME
	OFFTIME
)

//GetDifferences return an DeviceDifference bitmask which holds all the differences between the two shutters
//See const in model/comparer.go
func (s1 *Shutter) GetDifferences(s2 *Shutter) ModelDifference {
	result := NONE
	if s1.ClosePin != s2.ClosePin {
		result |= CLOSEPIN
	}
	if s1.OpenPin != s2.OpenPin {
		result |= OPENPIN
	}
	if s1.CompleteWayInSeconds != s2.CompleteWayInSeconds {
		result |= COMPLETEWAYINSECONDS
	}
	if s1.JobsEnabled != s2.JobsEnabled {
		result |= JOBSENABLED
	}
	if s1.OpenTime != s2.OpenTime {
		result |= OPENTIME
	}
	if s1.CloseTime != s2.CloseTime {
		result |= CLOSETIME
	}
	if s1.EmergencyEnabled != s2.EmergencyEnabled {
		result |= EMERGENCYENABLED
	}
	if s1.Disabled != s2.Disabled {
		result |= DISABLED
	}
	return result
}

//GetDifferences return an ModelDifferenceType bitmask which holds all the differences between the two shutters
//See const in model/comparer.go
func (l1 *Lighting) GetDifferences(l2 *Lighting) ModelDifference {
	result := NONE
	if l1.SwitchPin != l2.SwitchPin {
		result |= SWITCHPIN
	}
	if l1.JobsEnabled != l2.JobsEnabled {
		result |= JOBSENABLED
	}
	if l1.OnTime != l2.OnTime {
		result |= OPENTIME
	}
	if l1.OffTime != l2.OffTime {
		result |= CLOSETIME
	}
	if l1.EmergencyEnabled != l2.EmergencyEnabled {
		result |= EMERGENCYENABLED
	}
	if l1.Disabled != l2.Disabled {
		result |= DISABLED
	}
	return result
}
