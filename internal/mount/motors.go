package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

// SetMotorsEnabled enables or disables the mount's motor drivers (the TMC2209
// ENABLE line). Monturita keeps this line so the drivers can be idled; the
// firmware exposes it as the `motors` action and reports it back as
// `state.motors_enabled`.
func (m *Mount) SetMotorsEnabled(enabled bool) error {
	return m.action(mountproto.ActionRequest{Action: "motors", Enabled: &enabled})
}

// MotorsEnabled reports whether the motor drivers are enabled, from a fresh
// state snapshot.
func (m *Mount) MotorsEnabled() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return s.MotorsEnabled, nil
}
