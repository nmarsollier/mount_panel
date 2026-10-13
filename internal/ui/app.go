// Package ui exposes the mount to the Wails frontend. Every exported method on
// App is callable from JavaScript as window.go.main.App.<Method>().
package ui

import (
	"context"

	"github.com/nmarsollier/nrm_ha_go/internal/eaf"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// App is the Wails binding between the Alpine.js frontend and the mount.
type App struct {
	ctx context.Context
}

// NewApp builds the binding struct handed to wails.Run.
func NewApp() *App {
	return &App{}
}

// Startup is invoked once by Wails at launch.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Status returns the full mount snapshot for the UI.
func (a *App) Status() mount.Status {
	return mount.Get().Status()
}

// Stop halts everything (slew, tracking and manual motion).
func (a *App) Stop() error {
	return mount.Get().Stop()
}

// Home starts the home-finding sequence.
func (a *App) Home() error {
	return mount.Get().FindHome()
}

// Park moves the mount to the park position.
func (a *App) Park() error {
	return mount.Get().Park()
}

// Unpark takes the mount out of the parked state.
func (a *App) Unpark() error {
	return mount.Get().Unpark()
}

// Reset reboots the ESP.
func (a *App) Reset() error {
	return mount.Get().Reset()
}

// SetTracking sets the tracking mode: "none", "sidereal", "lunar" or "solar".
func (a *App) SetTracking(mode string) error {
	m := mount.Get()
	switch mode {
	case "sidereal":
		if err := m.SetTrackingRate(0); err != nil {
			return err
		}
		return m.SetTracking(true)
	case "lunar":
		if err := m.SetTrackingRate(1); err != nil {
			return err
		}
		return m.SetTracking(true)
	case "solar":
		if err := m.SetTrackingRate(2); err != nil {
			return err
		}
		return m.SetTracking(true)
	default:
		return m.SetTracking(false)
	}
}

// SetLimit fixes home or an axis limit from the current position, or clears.
func (a *App) SetLimit(action string) error {
	return mount.Get().SetLimits(action)
}

// MoveAxisSpeed sets both axes' manual rates at once (deg/s), for the joystick.
func (a *App) MoveAxisSpeed(raRate, decRate float64) error {
	return mount.Get().MoveAxisSpeed(raRate, decRate)
}

// MoveAxisBy moves one axis a relative number of degrees (manual slew).
func (a *App) MoveAxisBy(axis string, degrees float64, speed int) error {
	return mount.Get().MoveAxisBy(axis, degrees, speed)
}

// SlewTo slews to equatorial coordinates (RA hours, DEC degrees) at a speed.
func (a *App) SlewTo(ra, dec float64, speed int) error {
	return mount.Get().SlewTo(ra, dec, speed)
}

// SetTime re-syncs the mount clock to the host and clears any manual override,
// restoring the automatic time authority. The UI's "Update" button is the way
// back from a manual clock; an explicit Alpaca UTCDate set stays manual.
func (a *App) SetTime(_ string) error {
	return mount.Get().SyncTimeToHost()
}

// SetLocation updates the site latitude, longitude and elevation in one write.
func (a *App) SetLocation(lat, lon, elevation float64) error {
	return mount.Get().SetSite(lat, lon, int(elevation))
}

// EAFStatus returns the accessory (focuser/fan) snapshot for the UI.
func (a *App) EAFStatus() eaf.Status {
	return eaf.Get().Status()
}

// FocusMove moves the focuser a relative number of steps.
func (a *App) FocusMove(delta int64) error {
	return eaf.Get().MoveBy(delta)
}

// FocusStop stops the focuser.
func (a *App) FocusStop() error {
	return eaf.Get().Stop()
}

// FanSetPower switches the fan's 12 V relay.
func (a *App) FanSetPower(on bool) error {
	return eaf.Get().SetFanPower(on)
}

// FanSetSpeed sets the fan PWM duty (0..100 %).
func (a *App) FanSetSpeed(pct int) error {
	return eaf.Get().SetFanSpeed(pct)
}
