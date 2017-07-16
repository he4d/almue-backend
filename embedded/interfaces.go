package embedded

//DeviceStateStore must be implemented by the store that supports methods for updating the states of the devices
type DeviceStateStore interface {
	UpdateLightingState(int64, string) error

	UpdateShutterState(int64, string) error

	UpdateShutterOpening(int64, int) error
}
