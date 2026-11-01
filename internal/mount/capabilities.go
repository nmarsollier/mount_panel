package mount

// MoveAxisSpeeds returns the discrete MoveAxis speed ladder (deg/s) the mount
// exposes, ordered slowest to fastest.  The joystick uses this list directly.
func (m *Mount) MoveAxisSpeeds() ([]float64, error) {
	c, err := m.fetchCaps()
	if err != nil {
		return nil, err
	}
	return c.Capabilities.MoveAxisSpeedsDps, nil
}

// MinMoveAxisSpeed returns the slowest MoveAxis speed in deg/s.
func (m *Mount) MinMoveAxisSpeed() (float64, error) {
	c, err := m.fetchCaps()
	if err != nil {
		return 0, err
	}
	if len(c.Capabilities.MoveAxisSpeedsDps) == 0 {
		return 0.005, nil
	}
	return c.Capabilities.MoveAxisSpeedsDps[0], nil
}

// MaxMoveAxisSpeed returns the fastest MoveAxis speed in deg/s.
func (m *Mount) MaxMoveAxisSpeed() (float64, error) {
	c, err := m.fetchCaps()
	if err != nil {
		return 0, err
	}
	if len(c.Capabilities.MoveAxisSpeedsDps) == 0 {
		return 1.0, nil
	}
	return c.Capabilities.MoveAxisSpeedsDps[len(c.Capabilities.MoveAxisSpeedsDps)-1], nil
}
