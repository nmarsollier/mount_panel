package alpaca_eaf

import "github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"

// Devices returns the Alpaca devices exposed by the accessory controller:
// Focuser, ObservingConditions and Switch.
func Devices() []server.DeviceInfo {
	return []server.DeviceInfo{
		{DeviceName: deviceName, DeviceType: "Focuser", DeviceNumber: 0, UniqueID: "nrm-eaf-focuser-001"},
		{DeviceName: deviceName, DeviceType: "ObservingConditions", DeviceNumber: 0, UniqueID: "nrm-eaf-observingconditions-001"},
		{DeviceName: deviceName, DeviceType: "Switch", DeviceNumber: 0, UniqueID: "nrm-eaf-switch-001"},
	}
}
