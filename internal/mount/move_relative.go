package mount

import "github.com/nmarsollier/nrm_ha_go/internal/proto"

// MoveAxisBy moves a single axis a relative number of degrees (positive or
// negative), using the requested speed profile (1..4). This is the manual
// "nudge" operation; the ESP validates the axis and speed.
func (m *Mount) MoveAxisBy(axis string, degrees float64, speed int) error {
	return m.action(proto.ActionRequest{
		Action:  "move",
		Axis:    axis,
		Degrees: &degrees,
		Speed:   &speed,
	})
}
