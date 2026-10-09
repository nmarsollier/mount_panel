package alpaca

// Alpaca server identity — app-level constants, mirroring the firmware's
// rest_alpaca_internal.h values.
const (
	serverName        = "NRM-HA"
	serverDescription = "NRM-HA — Ecuatorial Mount Controller"
	driverInfo        = "NRM-HA Alpaca Driver v1.0"
	driverVersion     = "1.0.0"
	interfaceVersion  = 4

	deviceName   = "NRM-HA"
	deviceType   = "Telescope"
	deviceNumber = 0
	uniqueID     = "nrm-ha-telescope-001"

	apertureArea     = 0.0113 // m^2 (120mm)
	apertureDiameter = 0.120  // m (120mm)
	focalLength      = 0.900  // m (900mm)
)

// Alpaca error numbers (ASCOM standard).
const (
	errNotImplemented = 1024
	errInvalidValue   = 1025
	errNotConnected   = 1027 // ASCOM "Not Connected" (0x403): serial link down/unresponsive
)
