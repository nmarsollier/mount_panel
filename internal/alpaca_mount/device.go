package alpaca_mount

import (
	"github.com/nmarsollier/mount_panel/internal/alpaca_eaf"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
)

// configuredDevices returns the devices exposed by the server: the telescope
// plus the accessory controller's devices (focuser, observing conditions,
// switch).
func configuredDevices() []server.DeviceInfo {
	devices := []server.DeviceInfo{
		{DeviceName: deviceName, DeviceType: deviceType, DeviceNumber: deviceNumber, UniqueID: uniqueID},
	}
	return append(devices, alpaca_eaf.Devices()...)
}
