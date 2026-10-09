package mount

import (
	"errors"

	"github.com/nmarsollier/nrm_ha_go/internal/proto"
)

var guideDirections = [...]string{"north", "south", "east", "west"}

// PulseGuide executes a guide pulse in the given direction (0 North, 1 South,
// 2 East, 3 West) for the given duration in milliseconds.
func (m *Mount) PulseGuide(direction int, durationMs int) error {
	if direction < 0 || direction > 3 {
		return errors.New("Invalid guide direction")
	}
	d := uint32(durationMs)
	return m.action(proto.ActionRequest{
		Action:     "guide",
		Direction:  guideDirections[direction],
		DurationMs: &d,
	})
}

// SetGuideRateRA stores the RA guide rate in deg/s.
func (m *Mount) SetGuideRateRA(rate float64) error {
	return m.configWrite(proto.ConfigRequest{GuideRateRa: &rate})
}

// SetGuideRateDec stores the DEC guide rate in deg/s.
func (m *Mount) SetGuideRateDec(rate float64) error {
	return m.configWrite(proto.ConfigRequest{GuideRateDec: &rate})
}
