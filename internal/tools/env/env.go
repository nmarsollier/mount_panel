package env

import (
	"os"
	"strconv"
)

// Configuration holds the app settings, loaded from environment variables.
type Configuration struct {
	Port           int    `json:"port"`           // Alpaca REST port
	SerialPort     string `json:"serialPort"`     // mount CDC-ACM serial port (empty = auto-detect)
	EAFSerialPort  string `json:"eafSerialPort"`  // focuser CDC-ACM serial port (empty = auto-detect)
	FanSerialPort  string `json:"fanSerialPort"`  // fan CDC-ACM serial port (empty = auto-detect)
	TempSerialPort string `json:"tempSerialPort"` // temp/vibration CDC-ACM serial port (empty = auto-detect)
}

var config *Configuration

// Get returns the singleton Configuration, loading it on first use.
var Get = func() *Configuration {
	if config == nil {
		config = load()
	}
	return config
}

func load() *Configuration {
	port := 11111
	if v, err := strconv.Atoi(os.Getenv("PORT")); err == nil && v != 0 {
		port = v
	}

	return &Configuration{
		Port:           port,
		SerialPort:     os.Getenv("SERIAL_PORT"),
		EAFSerialPort:  os.Getenv("EAF_SERIAL_PORT"),
		FanSerialPort:  os.Getenv("FAN_SERIAL_PORT"),
		TempSerialPort: os.Getenv("TEMP_SERIAL_PORT"),
	}
}
