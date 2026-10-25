package eaf

import "github.com/nmarsollier/mount_panel/internal/tools/logx"

// Status is the UI snapshot of the accessory controller: connection state plus
// the focuser, fan, environment and vibration readings, each from a fresh
// serial read.
type Status struct {
	Detected     bool            `json:"detected"`
	BootID       string          `json:"boot_id"`
	Power        bool            `json:"power"`
	Health       string          `json:"health"`
	SensorHealth string          `json:"sensor_health"`
	Focuser      FocuserView     `json:"focuser"`
	Fan          FanView         `json:"fan"`
	Environment  EnvironmentView `json:"environment"`
	Vibration    VibrationView   `json:"vibration"`
}

// FocuserView is the focuser subset shown in the UI.
type FocuserView struct {
	CounterSteps int64  `json:"counter_steps"`
	IsMoving     bool   `json:"is_moving"`
	OpStatus     string `json:"op_status"`
	DriverOK     bool   `json:"driver_ok"`
	Enabled      bool   `json:"enabled"`
}

// FanView is the fan subset shown in the UI.
type FanView struct {
	Error       *string `json:"error"`
	Mode        string  `json:"mode"`
	Power       bool    `json:"power"`
	Speed       int     `json:"speed"`
	ManualSpeed int     `json:"manual_speed"`
	PowerOK     bool    `json:"power_ok"`
	RPM         *int    `json:"rpm"`
}

// EnvironmentView is the environment subset shown in the UI. A nil field means
// the sensor is absent or faulted; the frontend renders it as "—".
type EnvironmentView struct {
	Error               *string  `json:"error"`
	AmbientC            *float64 `json:"ambient_c"`
	HumidityPct         *float64 `json:"humidity_pct"`
	MirrorC             *float64 `json:"mirror_c"`
	DewPointC           *float64 `json:"dew_point_c"`
	MirrorAmbientDeltaC *float64 `json:"mirror_ambient_delta_c"`
	MirrorDewMarginC    *float64 `json:"mirror_dew_margin_c"`
}

// VibrationView is the vibration subset shown in the UI.
type VibrationView struct {
	Error    *string `json:"error"`
	RmsTotal float64 `json:"rms_total"`
}

// Status aggregates device, focuser, fan, environment and vibration readings.
// When the controller is not connected it returns a zero snapshot with
// Detected=false, so the UI can render the accessory controls disabled.
func (e *EAF) Status() Status {
	st := Status{Detected: false, Health: "unknown"}

	c, err := e.getConn()
	if err != nil {
		return st
	}
	st.Detected = true
	st.BootID = c.BootID()

	if dev, err := c.StateDevice(); err != nil {
		logx.Warn(eafCtx(), "state/device: "+err.Error())
	} else {
		st.Power = dev.PowerExternalPresent
		st.Health = dev.Health
		st.SensorHealth = dev.SensorHealth
	}
	if fs, err := c.StateFocuser(); err != nil {
		logx.Warn(eafCtx(), "state/focuser: "+err.Error())
	} else {
		st.Focuser = FocuserView{
			CounterSteps: fs.CounterSteps,
			IsMoving:     fs.IsMoving,
			OpStatus:     fs.OpStatus,
			DriverOK:     fs.Driver.UartOK,
			Enabled:      fs.Enabled,
		}
	}
	if fan, err := c.StateFan(); err != nil {
		logx.Warn(eafCtx(), "state/fan: "+err.Error())
	} else {
		st.Fan = FanView{
			Error:       fan.Error,
			Mode:        fan.Mode,
			Power:       fan.Power,
			Speed:       fan.Speed,
			ManualSpeed: fan.ManualSpeed,
			PowerOK:     fan.PowerAvailable,
			RPM:         fan.RPM,
		}
	}
	if env, err := c.StateEnvironment(); err != nil {
		logx.Warn(eafCtx(), "state/environment: "+err.Error())
	} else {
		st.Environment = EnvironmentView{
			Error:               env.Error,
			AmbientC:            env.AmbientC,
			HumidityPct:         env.HumidityPct,
			MirrorC:             env.MirrorC,
			DewPointC:           env.DewPointC,
			MirrorAmbientDeltaC: env.MirrorAmbientDeltaC,
			MirrorDewMarginC:    env.MirrorDewMarginC,
		}
	}
	if vib, err := c.StateVibration(); err != nil {
		logx.Warn(eafCtx(), "state/vibration: "+err.Error())
	} else {
		st.Vibration = VibrationView{
			Error:    vib.Error,
			RmsTotal: vib.RmsTotal,
		}
	}

	return st
}
