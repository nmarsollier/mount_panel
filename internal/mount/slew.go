package mount

import (
	"errors"
	"math"

	"github.com/nmarsollier/nrm_ha_go/internal/mountproto"
)

// SlewToCoordinates starts an asynchronous slew at speed profile 4 (Alpaca).
func (m *Mount) SlewToCoordinates(ra, dec float64) error {
	return m.SlewTo(ra, dec, 4)
}

// SlewTo starts an asynchronous slew to the given equatorial coordinates
// (RA in hours, DEC in degrees) at the requested speed profile (1..4).
func (m *Mount) SlewTo(ra, dec float64, speed int) error {
	if !isFinite(ra) || ra < 0 || ra >= 24 {
		return errors.New("RightAscension out of range [0,24)")
	}
	if !isFinite(dec) || dec < -90 || dec > 90 {
		return errors.New("Declination out of range [-90,90]")
	}
	return m.action(mountproto.ActionRequest{Action: "goto", Ra: &ra, Dec: &dec, Speed: &speed})
}

// SlewToTarget slews to the previously stored target coordinates.
func (m *Mount) SlewToTarget() error {
	m.mu.RLock()
	ra, dec := m.targetRa, m.targetDec
	m.mu.RUnlock()

	if ra == 0 && dec == 0 {
		return errors.New("Target not set")
	}
	speed := 2
	return m.action(mountproto.ActionRequest{Action: "goto", Ra: &ra, Dec: &dec, Speed: &speed})
}

// AbortSlew stops the in-progress slew while keeping an active tracking mode.
func (m *Mount) AbortSlew() error {
	s, err := m.fetchState()
	saved := ""
	if err == nil {
		saved = s.Tracking
	}

	if err := m.action(mountproto.ActionRequest{Action: "stop", Scope: "all"}); err != nil {
		return err
	}
	if saved != "" && saved != "none" {
		return m.control(mountproto.ControlRequest{Tracking: saved})
	}
	return nil
}

// TargetRA returns the stored target right ascension in hours.
func (m *Mount) TargetRA() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.targetRa
}

// TargetDec returns the stored target declination in degrees.
func (m *Mount) TargetDec() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.targetDec
}

// SetTargetRA stores the target right ascension (hours).
func (m *Mount) SetTargetRA(ra float64) error {
	if !isFinite(ra) || ra < 0 || ra >= 24 {
		return errors.New("TargetRightAscension out of range [0,24)")
	}
	m.mu.Lock()
	m.targetRa = ra
	m.mu.Unlock()
	return nil
}

// SetTargetDec stores the target declination (degrees).
func (m *Mount) SetTargetDec(dec float64) error {
	if !isFinite(dec) || dec < -90 || dec > 90 {
		return errors.New("TargetDeclination out of range [-90,90]")
	}
	m.mu.Lock()
	m.targetDec = dec
	m.mu.Unlock()
	return nil
}

func isFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}
