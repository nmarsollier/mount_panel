package proto

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.bug.st/serial"
)

// handshakeTimeout is how long a candidate port gets to answer the
// capabilities handshake during detection.
const handshakeTimeout = 1500 * time.Millisecond

var errNoMount = errors.New("no NRM-HA mount found on any serial port")

// Detect probes the serial ports and returns a client connected to the NRM-HA
// mount. Each candidate is verified with a capabilities handshake so a random
// serial device is never mistaken for the mount.
func Detect() (*Client, error) {
	ports, err := candidatePorts()
	if err != nil {
		return nil, err
	}

	for _, name := range ports {
		p, err := openPort(name)
		if err != nil {
			continue
		}

		c := newClient(p, name, handshakeTimeout)
		caps, err := c.Capabilities()
		if err == nil && isMount(caps) {
			c.timeout = defaultTimeout
			return c, nil
		}
		c.Close()
	}
	return nil, errNoMount
}

// candidatePorts returns the serial port names to probe, CDC-likely first.
func candidatePorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("enumerate serial ports: %w", err)
	}

	var cdc, other []string
	for _, name := range ports {
		if isCdcName(name) {
			cdc = append(cdc, name)
		} else {
			other = append(other, name)
		}
	}
	// /dev/cu.* sorts before /dev/tty.* so the outgoing node is tried first.
	sort.Strings(cdc)
	sort.Strings(other)
	return append(cdc, other...), nil
}

// isCdcName reports whether a port name looks like a USB CDC-ACM device.
func isCdcName(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "usbmodem") || // macOS
		strings.Contains(n, "ttyacm") || // Linux
		strings.Contains(n, "ttyusb") || // Linux / USB serial adapters
		strings.HasPrefix(n, "com") // Windows
}

// isMount reports whether a capabilities response identifies the mount: a
// valid handshake with a name and the ra/dec axes, which only the NRM-HA
// firmware produces.
func isMount(caps *CapabilitiesResponse) bool {
	return caps != nil && caps.Name != "" && caps.Proto >= 1 && len(caps.Axes) >= 2
}
