package model

//DifferenceType represents a bitmask with various flags which indicates the
//differences between two models
type DifferenceType uint16

const (
	NONE DifferenceType = 1 << iota
	EMERGENCYENABLED
	DISABLED
	JOBSENABLED
	OPENPIN
	CLOSEPIN
	COMPLETEWAYINSECONDS
	OPENTIME
	CLOSETIME
	SWITCHPIN
	ONTIME
	OFFTIME
)

//HasFlag checks if a ModelDifference bitmask has a specified flag
func (bitmask DifferenceType) HasFlag(flag DifferenceType) bool { return bitmask&flag != 0 }

//GetDifferences return an DeviceDifference bitmask which holds all the differences between the two shutters
//See const in model/comparer.go
func (s1 *Shutter) GetDifferences(s2 *Shutter) DifferenceType {
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
func (l1 *Lighting) GetDifferences(l2 *Lighting) DifferenceType {
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
