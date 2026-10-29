package proto

import "sync"

// DetectMu serialises device-detection scans across the mount and accessory
// clients. Detection probes every candidate port and claims each one
// atomically; running the scans concurrently lets one detector's probe of a
// port it does not match hold the port just long enough for another detector to
// skip it and miss its own device. Serialising the scans keeps the detectors
// from racing each other over the shared port-ownership map.
var DetectMu sync.Mutex

// portOwners tracks which serial ports are held by a client. The mount and
// accessory clients use separate ports, but both auto-detect by probing every
// candidate; the ownership map keeps one detector from opening (and corrupting)
// a port the other already owns or is probing. Detection runs concurrently:
// each candidate is claimed atomically before it is opened, and released again
// if the probe does not match the expected device.
var portOwners sync.Map // port name -> struct{}{}

// ClaimPort marks a serial port as owned by a client. It reports false when
// another client already owns it, so the other detector skips it.
func ClaimPort(name string) bool {
	_, loaded := portOwners.LoadOrStore(name, struct{}{})
	return !loaded
}

// ReleasePort frees a serial port so it can be detected again.
func ReleasePort(name string) {
	portOwners.Delete(name)
}
