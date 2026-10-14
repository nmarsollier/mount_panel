package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

// MoveAxisBy moves a single axis a relative number of degrees (positive or
// negative), using the requested speed profile (1..4). This is the manual
// "nudge" operation; the ESP validates the axis and speed.
func (m *Mount) MoveAxisBy(axis string, degrees float64, speed int) error {
	return m.action(mountproto.ActionRequest{
		Action:  "move",
		Axis:    axis,
		Degrees: &degrees,
		Speed:   &speed,
	})
}
