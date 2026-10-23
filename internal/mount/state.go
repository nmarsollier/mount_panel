package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

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

// pierUnknown is the ASCOM PierSide value for a fork mount: it has no side of
// pier, so SideOfPier and DestinationSideOfPier return it.
const pierUnknown = -1

// pierSideState returns the ASCOM PierSide (0 pierEast / 1 pierWest) parsed
// from the protocol pier_side string, for a German equatorial.
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

// SideOfPier returns pierUnknown (-1) for a fork, or 0 (pierEast) / 1
// (pierWest) for a German equatorial. The state is still fetched so a
// disconnected mount reports Not Connected.
func (m *Mount) SideOfPier() (int, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	if !m.MeridianFlip() {
		return pierUnknown, nil
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
