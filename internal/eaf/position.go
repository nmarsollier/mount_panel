package eaf

import "errors"

// Focuser limits for the Alpaca absolute position (configured in Go; the
// firmware enforces its own max_move_steps separately).
const (
	maxStep      = 1000000 // virtual position ceiling (steps)
	maxIncrement = 100000  // max per-move delta (steps)
)

var (
	// ErrNoReference is returned when an absolute move is requested but the
	// firmware's position reference is not established (counter_valid=false).
	ErrNoReference     = errors.New("focuser position reference not established")
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

// Position returns the last known absolute position. It mirrors the firmware's
// counter, so it is returned even when no reference is established (callers
// check PositionValid for that); it errors only when disconnected.
func (e *EAF) Position() (int32, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.conn == nil {
		return 0, errNotConnected
	}
	return e.position, nil
}

// PositionValid reports whether the firmware's position reference is current.
func (e *EAF) PositionValid() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.conn != nil && e.positionValid
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
	if !e.PositionValid() {
		return ErrNoReference
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

// SetPosition declares the focuser's current absolute position and marks the
// reference valid. It does not move the motor.
func (e *EAF) SetPosition(position int32) error {
	if position < 0 || position > maxStep {
		return errPositionRange
	}
	c, err := e.getConn()
	if err != nil {
		return err
	}
	resp, err := c.FocuserSetPosition(int64(position))
	if err != nil {
		return err
	}
	e.mu.Lock()
	e.position = int32(resp.CounterSteps)
	e.positionValid = resp.CounterValid
	e.mu.Unlock()
	return nil
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
