package mount

import (
	"errors"

	"github.com/nmarsollier/nrm_ha_go/internal/mountproto"
)

// rateToMode maps an ASCOM DriveRates value to the protocol tracking mode.
func rateToMode(rate int) string {
	switch rate {
	case 1:
		return "lunar"
	case 2:
		return "solar"
	default:
		return "sidereal"
	}
}

// SetTracking enables/disables tracking, keeping the currently selected rate.
func (m *Mount) SetTracking(enabled bool) error {
	m.mu.RLock()
	rate := m.selectedTrackingRate
	m.mu.RUnlock()

	mode := "none"
	if enabled {
		mode = rateToMode(rate)
	}
	return m.control(mountproto.ControlRequest{Tracking: mode})
}

// SelectedTrackingRate returns the last-set DriveRates value (the setpoint,
// not the active mode).
func (m *Mount) SelectedTrackingRate() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.selectedTrackingRate
}

// SetTrackingRate stores a DriveRates value and, if tracking is active, applies
// it immediately. It does not start tracking on its own.
func (m *Mount) SetTrackingRate(rate int) error {
	if rate < 0 || rate > 2 {
		return errors.New("Tracking rate out of range [0..2]")
	}

	m.mu.Lock()
	m.selectedTrackingRate = rate
	m.mu.Unlock()

	if trk, err := m.Tracking(); err == nil && trk {
		return m.control(mountproto.ControlRequest{Tracking: rateToMode(rate)})
	}
	return nil
}
