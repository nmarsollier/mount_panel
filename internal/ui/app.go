// Package ui exposes the mount to the Wails frontend. Every exported method on
// App is callable from JavaScript as window.go.main.App.<Method>().
package ui

import (
	"context"

	"github.com/nmarsollier/mount_panel/internal/eaf"
	"github.com/nmarsollier/mount_panel/internal/fan"
	"github.com/nmarsollier/mount_panel/internal/mount"
	"github.com/nmarsollier/mount_panel/internal/temp"
)

// App is the Wails binding between the Alpine.js frontend and the mount.
type App struct{}

// NewApp builds the binding struct handed to wails.Run.
func NewApp() *App {
	return &App{}
}

// Startup is invoked once by Wails at launch.
func (*App) Startup(context.Context) {
}

// Status returns the full mount snapshot for the UI.
func (a *App) Status() mount.Status {
	return mount.Get().Status()
}

// ConnectMount makes a single manual connection attempt to the mount.
func (a *App) ConnectMount() error {
	return mount.Get().Connect()
}

// ConnectEAF makes a single manual connection attempt to the focuser.
func (a *App) ConnectEAF() error {
	return eaf.Get().Connect()
}

// ConnectFan makes a single manual connection attempt to the fan controller.
func (a *App) ConnectFan() error {
	return fan.Get().Connect()
}

// ConnectTemp makes a single manual connection attempt to the temp module.
func (a *App) ConnectTemp() error {
	return temp.Get().Connect()
}

// DisconnectMount closes the serial link to the mount.
func (a *App) DisconnectMount() {
	mount.Get().Disconnect()
}

// DisconnectEAF closes the serial link to the focuser.
func (a *App) DisconnectEAF() {
	eaf.Get().Disconnect()
}

// DisconnectFan closes the serial link to the fan controller.
func (a *App) DisconnectFan() {
	fan.Get().Disconnect()
}

// DisconnectTemp closes the serial link to the temp module.
func (a *App) DisconnectTemp() {
	temp.Get().Disconnect()
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

// FocusMove moves the focuser a relative number of steps.
func (a *App) FocusMove(delta int64) error {
	return eaf.Get().MoveBy(delta)
}

// FocusStop stops the focuser.
func (a *App) FocusStop() error {
	return eaf.Get().Stop()
}

// FocusSetEnabled enables or disables the focuser driver.
func (a *App) FocusSetEnabled(enabled bool) error {
	return eaf.Get().SetEnabled(enabled)
}

// FanSetMode sets the fan mode: "off", "auto" or "manual".
func (a *App) FanSetMode(mode string) error {
	return fan.Get().SetFanMode(mode)
}

// FanSetSpeed sets the manual fan speed (1..100).
func (a *App) FanSetSpeed(pct int) error {
	return fan.Get().SetFanSpeed(pct)
}

// eafStatus is the aggregate accessory snapshot returned to the frontend: the
// focuser comes from the eaf, the fan and mirror from the fan controller, and
// the ambient/vibration readings from the temp module.
type eafStatus struct {
	Detected     bool            `json:"detected"`
	BootID       string          `json:"boot_id"`
	Power        bool            `json:"power"`
	Health       string          `json:"health"`
	SensorHealth string          `json:"sensor_health"`
	FanDetected  bool            `json:"fan_detected"`
	TempDetected bool            `json:"temp_detected"`
	Focuser      eaf.FocuserView `json:"focuser"`
	Fan          fan.FanView     `json:"fan"`
	Environment  environmentView `json:"environment"`
	Vibration    vibrationView   `json:"vibration"`
}

// environmentView merges the ambient readings (temp) with the mirror reading
// (fan) under the shape the frontend expects.
type environmentView struct {
	Error               *string  `json:"error"`
	AmbientC            *float64 `json:"ambient_c"`
	HumidityPct         *float64 `json:"humidity_pct"`
	MirrorC             *float64 `json:"mirror_c"`
	DewPointC           *float64 `json:"dew_point_c"`
	MirrorAmbientDeltaC *float64 `json:"mirror_ambient_delta_c"`
	MirrorDewMarginC    *float64 `json:"mirror_dew_margin_c"`
}

// vibrationView is the vibration subset shown in the UI.
type vibrationView struct {
	Error    *string `json:"error"`
	RmsTotal float64 `json:"rms_total"`
}

// EAFStatus returns the aggregate accessory snapshot for the UI, merging the
// three accessory controllers into the single shape the frontend renders.
func (a *App) EAFStatus() eafStatus {
	st := eafStatus{Health: "unknown", SensorHealth: "unknown"}

	fs := eaf.Get().Status()
	st.Detected = fs.Detected
	st.BootID = fs.BootID
	st.Power = fs.Power
	st.Health = fs.Health
	st.Focuser = fs.Focuser

	fst := fan.Get().Status()
	st.FanDetected = fst.Detected
	st.Fan = fst.Fan
	mirror := fst.Mirror.MirrorC
	mirrorErr := fst.Mirror.Error
	if !fst.Detected {
		mirrorErr = strptr("fan not connected")
	}

	tst := temp.Get().Status()
	st.TempDetected = tst.Detected
	st.SensorHealth = tst.SensorHealth
	ambientErr := tst.Environment.Error
	if !tst.Detected {
		ambientErr = strptr("temp not connected")
	}
	st.Vibration = vibrationView{Error: tst.Vibration.Error, RmsTotal: tst.Vibration.RmsTotal}
	if !tst.Detected && st.Vibration.Error == nil {
		st.Vibration.Error = strptr("temp not connected")
	}

	st.Environment = environmentView{
		Error:       mergeErr(ambientErr, mirrorErr),
		AmbientC:    tst.Environment.AmbientC,
		HumidityPct: tst.Environment.HumidityPct,
		MirrorC:     mirror,
		DewPointC:   tst.Environment.DewPointC,
	}

	return st
}

func strptr(s string) *string { return &s }

func mergeErr(a, b *string) *string {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	joined := *a + " · " + *b
	return &joined
}
