package eaf

import "github.com/nmarsollier/mount_panel/internal/tools/logx"

// FocuserView is the focuser subset shown in the UI.
type FocuserView struct {
	CounterSteps int64  `json:"counter_steps"`
	IsMoving     bool   `json:"is_moving"`
	OpStatus     string `json:"op_status"`
	DriverOK     bool   `json:"driver_ok"`
	Enabled      bool   `json:"enabled"`
}

// Status is the focuser snapshot: connection state plus the focuser reading.
type Status struct {
	Detected bool        `json:"detected"`
	BootID   string      `json:"boot_id"`
	Power    bool        `json:"power"`
	Health   string      `json:"health"`
	Focuser  FocuserView `json:"focuser"`
}

// Status aggregates the device and focuser readings. When the controller is
// not connected it returns a zero snapshot with Detected=false.
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

	return st
}
