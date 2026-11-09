package amb

import (
	"github.com/nmarsollier/mount_panel/internal/events"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// EnvironmentView is the ambient environment subset shown in the UI. The mirror
// fields are nil for this module (the mirror sensor lives in the fan).
type EnvironmentView struct {
	Error       *string  `json:"error"`
	AmbientC    *float64 `json:"ambient_c"`
	HumidityPct *float64 `json:"humidity_pct"`
	DewPointC   *float64 `json:"dew_point_c"`
}

// VibrationView is the vibration subset shown in the UI: acceleration RMS/peak
// (milli-g) and velocity RMS/peak (mm/s, converted from the firmware's µm/s),
// plus the ISO 10816 zone. A field is nil when the accelerometer is absent; the
// frontend renders it as "—".
type VibrationView struct {
	Error          *string  `json:"error"`
	RmsTotal       *float64 `json:"rms_total"`
	AccelPeakTotal *float64 `json:"accel_peak_total"`
	VelRmsTotal    *float64 `json:"vel_rms_total"`
	VelPeakTotal   *float64 `json:"vel_peak_total"`
	IsoZone        *string  `json:"iso_zone"`
}

// OrientationView is the tilt subset shown in the UI: pitch/roll tilt angles
// (degrees) and the gravity vector (g).
type OrientationView struct {
	Error    *string  `json:"error"`
	PitchDeg *float64 `json:"pitch_deg"`
	RollDeg  *float64 `json:"roll_deg"`
	AxG      *float64 `json:"ax_g"`
	AyG      *float64 `json:"ay_g"`
	AzG      *float64 `json:"az_g"`
}

// Status is the ambient/vibration snapshot: connection state plus the
// ambient, vibration and orientation readings.
type Status struct {
	Detected     bool            `json:"detected"`
	BootID       string          `json:"boot_id"`
	Health       string          `json:"health"`
	SensorHealth string          `json:"sensor_health"`
	Environment  EnvironmentView `json:"environment"`
	Vibration    VibrationView   `json:"vibration"`
	Orientation  OrientationView `json:"orientation"`
}

// Status aggregates the device, environment and vibration readings. When the
// module is not connected it returns a zero snapshot with Detected=false.
func (t *Amb) Status() Status {
	st := Status{Detected: false, Health: "unknown"}

	c, err := t.getConn()
	if err != nil {
		return st
	}
	st.Detected = true
	st.BootID = c.BootID()

	if dev, err := c.StateDevice(); err != nil {
		logx.Warn(ambCtx(), "state/device: "+err.Error())
	} else {
		st.Health = dev.Health
		st.SensorHealth = dev.SensorHealth
	}
	if env, err := c.StateEnvironment(); err != nil {
		logx.Warn(ambCtx(), "state/environment: "+err.Error())
	} else {
		st.Environment = EnvironmentView{
			Error:       env.Error,
			AmbientC:    env.AmbientC,
			HumidityPct: env.HumidityPct,
			DewPointC:   env.DewPointC,
		}
	}
	if vib, err := c.StateVibration(); err != nil {
		logx.Warn(ambCtx(), "state/vibration: "+err.Error())
	} else {
		st.Vibration = VibrationView{
			Error:          vib.Error,
			RmsTotal:       vib.RmsTotal,
			AccelPeakTotal: vib.AccelPeakTotal,
			VelRmsTotal:    umPerSToMmPerS(vib.VelRmsTotal),
			VelPeakTotal:   umPerSToMmPerS(vib.VelPeakTotal),
			IsoZone:        vib.IsoZone,
		}
	}
	if ori, err := c.StateOrientation(); err != nil {
		logx.Warn(ambCtx(), "state/orientation: "+err.Error())
	} else {
		st.Orientation = OrientationView{
			Error:    ori.Error,
			PitchDeg: ori.PitchDeg,
			RollDeg:  ori.RollDeg,
			AxG:      ori.AxG,
			AyG:      ori.AyG,
			AzG:      ori.AzG,
		}
	}

	if st.Detected {
		events.Record(events.Event{
			Type:    events.TypeSensorReading,
			Source:  events.SourceAmb,
			Subject: "status",
			Details: events.SensorAmb{
				Detected:       st.Detected,
				BootID:         st.BootID,
				Health:         st.Health,
				SensorHealth:   st.SensorHealth,
				EnvError:       st.Environment.Error,
				AmbientC:       st.Environment.AmbientC,
				HumidityPct:    st.Environment.HumidityPct,
				DewPointC:      st.Environment.DewPointC,
				VibError:       st.Vibration.Error,
				RmsTotal:       st.Vibration.RmsTotal,
				AccelPeakTotal: st.Vibration.AccelPeakTotal,
				VelRmsTotal:    st.Vibration.VelRmsTotal,
				VelPeakTotal:   st.Vibration.VelPeakTotal,
				IsoZone:        st.Vibration.IsoZone,
				OriError:       st.Orientation.Error,
				PitchDeg:       st.Orientation.PitchDeg,
				RollDeg:        st.Orientation.RollDeg,
				AxG:            st.Orientation.AxG,
				AyG:            st.Orientation.AyG,
				AzG:            st.Orientation.AzG,
			},
		})
	}

	return st
}

// umPerSToMmPerS converts the firmware's µm/s velocity to mm/s for display. A
// nil value (sensor absent) is passed through.
func umPerSToMmPerS(v *float64) *float64 {
	if v == nil {
		return nil
	}
	m := *v / 1000.0
	return &m
}
