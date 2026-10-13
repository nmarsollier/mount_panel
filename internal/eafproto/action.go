package eafproto

import "context"

// FocuserMoveResponse is the accepted move result. ok=true means the move was
// accepted, not that it finished.
type FocuserMoveResponse struct {
	OpID          string `json:"op_id"`
	Status        string `json:"status"`
	IsMoving      bool   `json:"is_moving"`
	CounterSteps  int64  `json:"counter_steps"`
	TargetCounter int64  `json:"target_counter"`
}

// FocuserStopResponse is the result of stopping the focuser.
type FocuserStopResponse struct {
	IsMoving      bool   `json:"is_moving"`
	CounterSteps  int64  `json:"counter_steps"`
	CancelledOpID string `json:"cancelled_op_id"`
}

// FocuserMove sends a relative move of deltaSteps pulses, guarded by
// expectedCounter: the firmware rejects the move with COUNTER_CHANGED if the
// counter no longer matches (a local button moved the motor meanwhile).
func (c *Client) FocuserMove(deltaSteps int64, expectedCounter int64) (*FocuserMoveResponse, error) {
	out := struct {
		Data FocuserMoveResponse `json:"data"`
	}{}
	err := c.request(context.Background(), "action", "focuser", "move",
		map[string]any{"delta_steps": deltaSteps, "expected_counter": expectedCounter}, &out)
	if err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// FocuserStop stops the focuser and returns the final counter.
func (c *Client) FocuserStop() (*FocuserStopResponse, error) {
	out := struct {
		Data FocuserStopResponse `json:"data"`
	}{}
	if err := c.request(context.Background(), "action", "focuser", "stop", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}
