package model

type DeviceDifference uint16

const (
	NONE DeviceDifference = 1 << iota
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

// //GetModelDifferences return an ModelDifferenceType bitmask which holds all the differences between the two shutters
// //See const in model/comparer.go
// func (shutter *shutter) GetDeviceDifferences(shutterModel *model.Shutter) DeviceDifference {
// 	result := NONE
// 	if shutter.closePin.Number() != *shutterModel.ClosePin {
// 		result |= CLOSEPIN
// 	}
// 	if shutter.openPin.Number() != *shutterModel.OpenPin {
// 		result |= OPENPIN
// 	}
// 	if int(shutter.completeWayDuration.Seconds()) != *shutterModel.CompleteWayInSeconds {
// 		result |= COMPLETEWAYINSECONDS
// 	}
// 	if shutter.jobsEnabled != shutterModel.TimerEnabled {
// 		result |= JOBSENABLED
// 	}
// 	if shutter.OpenTime != shutterModel.OpenTime {
// 		result |= OPENTIME
// 	}
// 	if shutter.CloseTime != shutterModel.CloseTime {
// 		result |= CLOSETIME
// 	}
// 	if shutter.EmergencyEnabled != shutterModel.EmergencyEnabled {
// 		result |= EMERGENCYENABLED
// 	}
// 	if shutter.Disabled != shutterModel.Disabled {
// 		result |= DISABLED
// 	}
// 	return result
// }

// //GetModelDifferences return an ModelDifferenceType bitmask which holds all the differences between the two shutters
// //See const in model/comparer.go
// func (oldLighting *Lighting) GetModelDifferences(newLighting *Lighting) DeviceDifference {
// 	result := NONE
// 	if oldLighting.SwitchPin != newLighting.SwitchPin {
// 		result |= SWITCHPIN
// 	}
// 	if oldLighting.TimerEnabled != newLighting.TimerEnabled {
// 		result |= TIMERENABLED
// 	}
// 	if oldLighting.OnTime != newLighting.OnTime {
// 		result |= OPENTIME
// 	}
// 	if oldLighting.OffTime != newLighting.OffTime {
// 		result |= CLOSETIME
// 	}
// 	if oldLighting.EmergencyEnabled != newLighting.EmergencyEnabled {
// 		result |= EMERGENCYENABLED
// 	}
// 	if oldLighting.Disabled != newLighting.Disabled {
// 		result |= DISABLED
// 	}
// 	return result
// }
