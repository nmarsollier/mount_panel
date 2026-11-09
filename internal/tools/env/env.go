package env

import (
	"os"
	"path/filepath"
	"strconv"
)

// Configuration holds the app settings, loaded from environment variables.
type Configuration struct {
	Port          int    `json:"port"`          // Alpaca REST port
	Host          string `json:"host"`          // Alpaca REST bind address (default 127.0.0.1 = localhost only)
	SerialPort    string `json:"serialPort"`    // mount CDC-ACM serial port (empty = auto-detect)
	EAFSerialPort string `json:"eafSerialPort"` // focuser CDC-ACM serial port (empty = auto-detect)
	FanSerialPort string `json:"fanSerialPort"` // fan CDC-ACM serial port (empty = auto-detect)
	AmbSerialPort string `json:"ambSerialPort"` // ambient/vibration CDC-ACM serial port (empty = auto-detect)
	DBPath        string `json:"dbPath"`        // SQLite event-log file path (default: user config dir)
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
		Port:          port,
		Host:          envOr("ALPACA_HOST", "127.0.0.1"),
		SerialPort:    os.Getenv("SERIAL_PORT"),
		EAFSerialPort: os.Getenv("EAF_SERIAL_PORT"),
		FanSerialPort: os.Getenv("FAN_SERIAL_PORT"),
		AmbSerialPort: os.Getenv("AMB_SERIAL_PORT"),
		DBPath:        envOr("DB_PATH", defaultDBPath()),
	}
}

// envOr returns the environment variable value when set, or the fallback.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// defaultDBPath returns the event database location inside the user config
// directory, falling back to the current directory when it cannot be resolved.
func defaultDBPath() string {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		return "mount_panel.db"
	}
	return filepath.Join(dir, "mount_panel", "events.db")
}
