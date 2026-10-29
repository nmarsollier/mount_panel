package fan

import "github.com/nmarsollier/mount_panel/internal/tools/logx"

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

// MirrorView is the mirror-temperature subset shown in the UI. A nil field
// means the sensor is absent or faulted; the frontend renders it as "—".
type MirrorView struct {
	Error   *string  `json:"error"`
	MirrorC *float64 `json:"mirror_c"`
}

// Status is the fan snapshot: connection state plus the fan and mirror readings.
type Status struct {
	Detected bool       `json:"detected"`
	BootID   string     `json:"boot_id"`
	Power    bool       `json:"power"`
	Health   string     `json:"health"`
	Fan      FanView    `json:"fan"`
	Mirror   MirrorView `json:"mirror"`
}

// Status aggregates the device, fan and mirror readings. When the controller is
// not connected it returns a zero snapshot with Detected=false.
func (f *Fan) Status() Status {
	st := Status{Detected: false, Health: "unknown"}

	c, err := f.getConn()
	if err != nil {
		return st
	}
	st.Detected = true
	st.BootID = c.BootID()

	if dev, err := c.StateDevice(); err != nil {
		logx.Warn(fanCtx(), "state/device: "+err.Error())
	} else {
		st.Power = dev.PowerExternalPresent
		st.Health = dev.Health
	}
	if fan, err := c.StateFan(); err != nil {
		logx.Warn(fanCtx(), "state/fan: "+err.Error())
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
		logx.Warn(fanCtx(), "state/environment: "+err.Error())
	} else {
		st.Mirror = MirrorView{
			Error:   env.Error,
			MirrorC: env.MirrorC,
		}
	}

	return st
}
