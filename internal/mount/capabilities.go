package mount

// MinSlewSpeed returns the slowest MoveAxis/slew speed in deg/s.
func (m *Mount) MinSlewSpeed() float64 {
	c := m.fetchCaps()
	if c == nil || len(c.SlewSpeedsDps) == 0 {
		return 1.0
	}
	return c.SlewSpeedsDps[0]
}

// MaxSlewSpeed returns the fastest MoveAxis/slew speed in deg/s.
func (m *Mount) MaxSlewSpeed() float64 {
	c := m.fetchCaps()
	if c == nil || len(c.SlewSpeedsDps) == 0 {
		return 6.0
	}
	return c.SlewSpeedsDps[len(c.SlewSpeedsDps)-1]
}
