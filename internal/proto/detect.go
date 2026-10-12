package proto

import (
	"fmt"
	"sort"
	"strings"

	"go.bug.st/serial"
)

// CandidatePorts returns the USB CDC-ACM port names to probe, sorted. The
// mount and accessory controllers are CDC devices; non-CDC ports (Bluetooth,
// debug consoles) never match and only add a handshake timeout each, so they
// are excluded.
func CandidatePorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("enumerate serial ports: %w", err)
	}

	var cdc []string
	for _, name := range ports {
		if IsCdcName(name) {
			cdc = append(cdc, name)
		}
	}
	sort.Strings(cdc)
	return cdc, nil
}

// IsCdcName reports whether a port name looks like a USB CDC-ACM device.
func IsCdcName(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "usbmodem") || // macOS
		strings.Contains(n, "ttyacm") || // Linux
		strings.Contains(n, "ttyusb") || // Linux / USB serial adapters
		strings.HasPrefix(n, "com") // Windows
}
