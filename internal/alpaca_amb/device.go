package alpaca_amb

import "github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"

// Devices returns the Alpaca devices exposed by the ambient controller: a
// single ObservingConditions V2 device.
func Devices() []server.DeviceInfo {
	return []server.DeviceInfo{
		{DeviceName: deviceName, DeviceType: "ObservingConditions", DeviceNumber: 0, UniqueID: uniqueID},
	}
}
