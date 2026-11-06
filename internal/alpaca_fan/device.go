package alpaca_fan

import "github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"

// Devices returns the Alpaca devices exposed by the fan controller: a single
// Switch V3 device.
func Devices() []server.DeviceInfo {
	return []server.DeviceInfo{
		{DeviceName: deviceName, DeviceType: "Switch", DeviceNumber: 0, UniqueID: uniqueID},
	}
}
