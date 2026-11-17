package alpaca_eaf

import (
	"errors"

	"github.com/nmarsollier/mount_panel/internal/amb"
	"github.com/nmarsollier/mount_panel/internal/fan"
)

// The focuser has no temperature sensor of its own: the ASCOM Temperature it
// publishes for temperature-compensated focus is derived from the two sensors
// the mount does expose — the ambient SHT31 (amb module) and the primary-mirror
// DS18B20 (fan module). The published value is the average of both when both
// are available, the single available reading when only one is, and nothing
// when neither is.

// errNoTemperature reports that neither the ambient nor the mirror sensor
// produced a reading, so the focuser has no temperature to publish.
var errNoTemperature = errors.New("no temperature sensor available")

// reportedTemperature returns the focuser temperature in Celsius.
func reportedTemperature() (float64, error) {
	ambientC, ambientErr := amb.Get().AmbientTemperature()
	mirrorC, mirrorErr := fan.Get().MirrorTemperature()

	switch {
	case ambientErr == nil && mirrorErr == nil:
		return (ambientC + mirrorC) / 2, nil
	case ambientErr == nil:
		return ambientC, nil
	case mirrorErr == nil:
		return mirrorC, nil
	default:
		return 0, errNoTemperature
	}
}
