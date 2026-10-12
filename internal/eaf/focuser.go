package eaf

import (
	"errors"

	"github.com/nmarsollier/nrm_ha_go/internal/eafproto"
)

// maxMoveRetries bounds the COUNTER_CHANGED reconciliation loop, so a counter
// that keeps changing (a local button held down) does not spin forever.
const maxMoveRetries = 3

// errCounterChanged is returned when the focuser counter kept changing across
// the bounded retries, so the move could not be accepted safely.
var errCounterChanged = errors.New("focuser counter keeps changing (local control active?)")

// MoveBy moves the focuser a relative number of steps. The sign follows the
// firmware STEP counter (positive and negative are both valid). The move is
// guarded by the current counter, so a local button that moved the motor
// between the read and the command is detected (COUNTER_CHANGED) and retried.
func (e *EAF) MoveBy(deltaSteps int64) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}

	for attempt := 0; attempt < maxMoveRetries; attempt++ {
		fs, err := c.StateFocuser()
		if err != nil {
			return err
		}
		_, err = c.FocuserMove(deltaSteps, fs.CounterSteps)
		if err == nil {
			e.mu.Lock()
			e.position += int32(deltaSteps)
			e.mu.Unlock()
			return nil
		}
		if !isCounterChanged(err) {
			return err
		}
	}
	return errCounterChanged
}

// isCounterChanged reports whether the firmware rejected a move because the
// counter moved since it was read.
func isCounterChanged(err error) bool {
	var e *eafproto.Error
	return errors.As(err, &e) && e.Code == "COUNTER_CHANGED"
}

// Halt stops the focuser immediately.
func (e *EAF) Halt() error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	_, err = c.FocuserHalt()
	return err
}
