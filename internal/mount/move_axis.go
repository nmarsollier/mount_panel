package mount

import (
	"errors"

	"github.com/nmarsollier/nrm_ha_go/internal/proto"
)

// MoveAxis moves a single axis (0 = RA, 1 = DEC) continuously at the given
// rate in deg/s. The ESP clamps the rate to its safe ceiling.
func (m *Mount) MoveAxis(axis int, rate float64) error {
	req := proto.ControlRequest{}
	switch axis {
	case 0:
		req.ManualRaDps = &rate
	case 1:
		req.ManualDecDps = &rate
	default:
		return errors.New("Axis not supported")
	}
	return m.control(req)
}

// MoveAxisSpeed sets both axes' manual rates at once (deg/s); a zero rate
// stops that axis. Used by the joystick "hold to move" control.
func (m *Mount) MoveAxisSpeed(raRate, decRate float64) error {
	return m.control(proto.ControlRequest{
		ManualRaDps:  &raRate,
		ManualDecDps: &decRate,
	})
}
