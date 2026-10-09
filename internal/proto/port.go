package proto

import "go.bug.st/serial"

// openPort opens the given serial port at the CDC baud rate.
func openPort(name string) (serial.Port, error) {
	return serial.Open(name, &serial.Mode{BaudRate: 115200})
}
