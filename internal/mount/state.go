package mount

import "github.com/nmarsollier/nrm_ha_go/internal/mountproto"

// slewingState reports whether the state snapshot shows a slew in progress.
func slewingState(s *mountproto.StateResponse) bool {
	return s.State == "slewing"
}

// trackingState reports whether any automatic tracking mode is active.
func trackingState(s *mountproto.StateResponse) bool {
	if s.State == "error" {
		return false
	}
	return s.Tracking != "none"
}

// pierSideState returns 0 (pierEast) or 1 (pierWest).
func pierSideState(s *mountproto.StateResponse) int {
	if s.PierSide == "west" {
		return 1
	}
	return 0
}

// RightAscension returns the current right ascension in hours.
func (m *Mount) RightAscension() (float64, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	return s.Ra, nil
}

// Declination returns the current declination in degrees.
func (m *Mount) Declination() (float64, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	return s.Dec, nil
}

// SideOfPier returns 0 (pierEast) or 1 (pierWest).
func (m *Mount) SideOfPier() (int, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	return pierSideState(s), nil
}

// Slewing reports whether a slew is in progress.
func (m *Mount) Slewing() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return slewingState(s), nil
}

// Tracking reports whether any automatic tracking mode is active.
func (m *Mount) Tracking() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return trackingState(s), nil
}

// AtPark reports whether the mount is parked.
func (m *Mount) AtPark() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return s.AtPark, nil
}

// AtHome reports whether the mount is at the home position.
func (m *Mount) AtHome() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return s.AtHome, nil
}

// IsPulseGuiding reports whether a guide pulse is active.
func (m *Mount) IsPulseGuiding() (bool, error) {
	s, err := m.fetchState()
	if err != nil {
		return false, err
	}
	return s.Guiding, nil
}
