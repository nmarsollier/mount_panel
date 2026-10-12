package server

// DeviceInfo identifies one configured Alpaca device, as returned by the
// management configureddevices endpoint.
type DeviceInfo struct {
	DeviceName   string `json:"DeviceName"`
	DeviceType   string `json:"DeviceType"`
	DeviceNumber int    `json:"DeviceNumber"`
	UniqueID     string `json:"UniqueID"`
}
