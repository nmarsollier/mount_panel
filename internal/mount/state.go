package mount

// RightAscension returns the current right ascension in hours (0 if unknown).
func (m *Mount) RightAscension() float64 {
	s := m.fetchState()
	if s == nil {
		return 0
	}
	return s.Ra
}

// Declination returns the current declination in degrees (0 if unknown).
func (m *Mount) Declination() float64 {
	s := m.fetchState()
	if s == nil {
		return 0
	}
	return s.Dec
}

// SideOfPier returns 0 (pierEast) or 1 (pierWest).
func (m *Mount) SideOfPier() int {
	s := m.fetchState()
	if s == nil || s.PierSide == "west" {
		return 1
	}
	return 0
}

// Slewing reports whether a slew is in progress.
func (m *Mount) Slewing() bool {
	s := m.fetchState()
	return s != nil && s.State == "slewing"
}

// Tracking reports whether any automatic tracking mode is active.
func (m *Mount) Tracking() bool {
	s := m.fetchState()
	if s == nil {
		return false
	}
	if s.State == "error" {
		return false
	}
	return s.Tracking != "none"
}

// AtPark reports whether the mount is parked.
func (m *Mount) AtPark() bool {
	s := m.fetchState()
	return s != nil && s.AtPark
}

// AtHome reports whether the mount is at the home position.
func (m *Mount) AtHome() bool {
	s := m.fetchState()
	return s != nil && s.AtHome
}

// IsPulseGuiding reports whether a guide pulse is active.
func (m *Mount) IsPulseGuiding() bool {
	s := m.fetchState()
	return s != nil && s.Guiding
}
