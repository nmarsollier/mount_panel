package temp

import "github.com/nmarsollier/mount_panel/internal/tools/logx"

// EnvironmentView is the ambient environment subset shown in the UI. The mirror
// fields are nil for this module (the mirror sensor lives in the fan).
type EnvironmentView struct {
	Error       *string  `json:"error"`
	AmbientC    *float64 `json:"ambient_c"`
	HumidityPct *float64 `json:"humidity_pct"`
	DewPointC   *float64 `json:"dew_point_c"`
}

// VibrationView is the vibration subset shown in the UI.
type VibrationView struct {
	Error    *string `json:"error"`
	RmsTotal float64 `json:"rms_total"`
}

// Status is the temperature/vibration snapshot: connection state plus the
// ambient and vibration readings.
type Status struct {
	Detected     bool             `json:"detected"`
	BootID       string           `json:"boot_id"`
	Health       string           `json:"health"`
	SensorHealth string           `json:"sensor_health"`
	Environment  EnvironmentView  `json:"environment"`
	Vibration    VibrationView    `json:"vibration"`
}

// Status aggregates the device, environment and vibration readings. When the
// module is not connected it returns a zero snapshot with Detected=false.
func (t *Temp) Status() Status {
	st := Status{Detected: false, Health: "unknown"}

	c, err := t.getConn()
	if err != nil {
		return st
	}
	st.Detected = true
	st.BootID = c.BootID()

	if dev, err := c.StateDevice(); err != nil {
		logx.Warn(tempCtx(), "state/device: "+err.Error())
	} else {
		st.Health = dev.Health
		st.SensorHealth = dev.SensorHealth
	}
	if env, err := c.StateEnvironment(); err != nil {
		logx.Warn(tempCtx(), "state/environment: "+err.Error())
	} else {
		st.Environment = EnvironmentView{
			Error:       env.Error,
			AmbientC:    env.AmbientC,
			HumidityPct: env.HumidityPct,
			DewPointC:   env.DewPointC,
		}
	}
	if vib, err := c.StateVibration(); err != nil {
		logx.Warn(tempCtx(), "state/vibration: "+err.Error())
	} else {
		st.Vibration = VibrationView{
			Error:    vib.Error,
			RmsTotal: vib.RmsTotal,
		}
	}

	return st
}
