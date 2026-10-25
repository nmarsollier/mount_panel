package eaf

import "errors"

// Focuser limits for the Alpaca absolute position (configured in Go; the
// firmware enforces its own max_move_steps separately).
const (
	maxStep      = 1000000 // virtual position ceiling (steps)
	maxIncrement = 100000  // max per-move delta (steps)
)

var (
	errNoReference     = errors.New("focuser position reference not established")
	errPositionRange   = errors.New("focuser position out of range")
	errStepSizeUnknown = errors.New("focuser step size not calibrated")
)

// Connected reports the logical client link (Alpaca Connected).
func (e *EAF) Connected() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.connected
}

// SetConnected sets the logical client link. It never moves the focuser.
func (e *EAF) SetConnected(connected bool) {
	e.mu.Lock()
	e.connected = connected
	e.mu.Unlock()
}

// Position returns the virtual absolute position, or an error when no
// reference has been established.
func (e *EAF) Position() (int32, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !e.positionValid {
		return 0, errNoReference
	}
	return e.position, nil
}

// MaxStep returns the virtual position ceiling.
func (e *EAF) MaxStep() int32 { return maxStep }

// MaxIncrement returns the maximum per-move delta.
func (e *EAF) MaxIncrement() int32 { return maxIncrement }

// StepSize returns the calibrated step size in microns, or an error while it
// is unknown.
func (e *EAF) StepSize() (float64, error) {
	return 0, errStepSizeUnknown
}

// MoveTo moves the focuser to an absolute virtual position.
func (e *EAF) MoveTo(position int32) error {
	if position < 0 || position > maxStep {
		return errPositionRange
	}
	cur, err := e.Position()
	if err != nil {
		return err
	}
	delta := int64(position) - int64(cur)
	if delta == 0 {
		return nil
	}
	if delta > maxIncrement || delta < -maxIncrement {
		return errPositionRange
	}
	return e.MoveBy(delta)
}

// IsMoving reports whether the focuser is currently moving.
func (e *EAF) IsMoving() (bool, error) {
	c, err := e.getConn()
	if err != nil {
		return false, err
	}
	fs, err := c.StateFocuser()
	if err != nil {
		return false, err
	}
	return fs.IsMoving, nil
}

// Temperature returns the temperature in Celsius to report to Alpaca: the
// average of the ambient and mirror sensors when both are present, the only
// working sensor when one is absent, or ErrEnvironmentUnavailable when neither
// has a valid reading.
func (e *EAF) Temperature() (float64, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	ambient, mirror := env.AmbientC, env.MirrorC
	switch {
	case ambient != nil && mirror != nil:
		return (*ambient + *mirror) / 2, nil
	case ambient != nil:
		return *ambient, nil
	case mirror != nil:
		return *mirror, nil
	default:
		return 0, ErrEnvironmentUnavailable
	}
}
