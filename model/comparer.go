package model

//DifferenceType represents a bitmask with various flags which indicates the
//differences between two models
type DifferenceType uint16

const (
	// DIFFNONE identifies no differnces
	DIFFNONE DifferenceType = 1 << iota
	// DIFFEMERGENCYENABLED identifies different emergency enabled status
	DIFFEMERGENCYENABLED
	// DIFFDISABLED identifies different device enabled status
	DIFFDISABLED
	// DIFFJOBSENABLED identifies different jobs enabled status
	DIFFJOBSENABLED
	// DIFFOPENPIN identifies different open pin numbers
	DIFFOPENPIN
	// DIFFCLOSEPIN identifies different close pin
	DIFFCLOSEPIN
	// DIFFCOMPLETEWAYINSECONDS identifies different numbers complete way in secs
	DIFFCOMPLETEWAYINSECONDS
	// DIFFOPENTIME identifies different open time
	DIFFOPENTIME
	// DIFFCLOSETIME identifies different close time
	DIFFCLOSETIME
	// DIFFSWITCHPIN identifies different switch pin
	DIFFSWITCHPIN
	// DIFFONTIME identifies different on time
	DIFFONTIME
	// DIFFOFFTIME identifies different off time
	DIFFOFFTIME
)

//HasFlag checks if a ModelDifference bitmask has a specified flag
func (bitmask DifferenceType) HasFlag(flag DifferenceType) bool { return bitmask&flag != 0 }

//GetDifferences return an DeviceDifference bitmask which holds all the differences between the two shutters
//See const in model/comparer.go
func (s1 *Shutter) GetDifferences(s2 *Shutter) DifferenceType {
	result := DIFFNONE
	if *s1.ClosePin != *s2.ClosePin {
		result |= DIFFCLOSEPIN
	}
	if *s1.OpenPin != *s2.OpenPin {
		result |= DIFFOPENPIN
	}
	if *s1.CompleteWayInSeconds != *s2.CompleteWayInSeconds {
		result |= DIFFCOMPLETEWAYINSECONDS
	}
	if s1.JobsEnabled != s2.JobsEnabled {
		result |= DIFFJOBSENABLED
	}
	if s1.OpenTime != s2.OpenTime {
		result |= DIFFOPENTIME
	}
	if s1.CloseTime != s2.CloseTime {
		result |= DIFFCLOSETIME
	}
	if s1.EmergencyEnabled != s2.EmergencyEnabled {
		result |= DIFFEMERGENCYENABLED
	}
	if s1.Disabled != s2.Disabled {
		result |= DIFFDISABLED
	}
	return result
}

//GetDifferences return an ModelDifferenceType bitmask which holds all the differences between the two shutters
//See const in model/comparer.go
func (l1 *Lighting) GetDifferences(l2 *Lighting) DifferenceType {
	result := DIFFNONE
	if l1.SwitchPin != l2.SwitchPin {
		result |= DIFFSWITCHPIN
	}
	if l1.JobsEnabled != l2.JobsEnabled {
		result |= DIFFJOBSENABLED
	}
	if l1.OnTime != l2.OnTime {
		result |= DIFFOPENTIME
	}
	if l1.OffTime != l2.OffTime {
		result |= DIFFCLOSETIME
	}
	if l1.EmergencyEnabled != l2.EmergencyEnabled {
		result |= DIFFEMERGENCYENABLED
	}
	if l1.Disabled != l2.Disabled {
		result |= DIFFDISABLED
	}
	return result
}
