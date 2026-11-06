package alpaca_eaf

import "github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"

// Devices returns the Alpaca devices exposed by the focuser controller: a
// single Focuser V4 device.
func Devices() []server.DeviceInfo {
	return []server.DeviceInfo{
		{DeviceName: deviceName, DeviceType: deviceType, DeviceNumber: deviceNumber, UniqueID: uniqueID},
	}
}
