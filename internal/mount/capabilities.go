package mount

// MinSlewSpeed returns the slowest MoveAxis/slew speed in deg/s.
func (m *Mount) MinSlewSpeed() (float64, error) {
	c, err := m.fetchCaps()
	if err != nil {
		return 0, err
	}
	if len(c.SlewSpeedsDps) == 0 {
		return 1.0, nil
	}
	return c.SlewSpeedsDps[0], nil
}

// MaxSlewSpeed returns the fastest MoveAxis/slew speed in deg/s.
func (m *Mount) MaxSlewSpeed() (float64, error) {
	c, err := m.fetchCaps()
	if err != nil {
		return 0, err
	}
	if len(c.SlewSpeedsDps) == 0 {
		return 6.0, nil
	}
	return c.SlewSpeedsDps[len(c.SlewSpeedsDps)-1], nil
}
