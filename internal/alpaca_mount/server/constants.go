package server

// Alpaca error numbers (ASCOM standard).
const (
	ErrNotImplemented = 1024
	ErrInvalidValue   = 1025
	ErrNotConnected   = 1027 // ASCOM "Not Connected" (0x403)
)

// Alpaca server identity — shared across the device types served by this
// gateway.
const (
	ServerName        = "Mount Panel"
	ServerDescription = "Mount Panel — Mount + Accessories Controller"
	DriverInfo        = "Mount Panel Alpaca Driver v1.0"
	DriverVersion     = "1.0.0"
	InterfaceVersion  = 4
)
