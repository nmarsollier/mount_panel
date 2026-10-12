package mountproto

import "context"

// ActionRequest is an operation with a beginning and an end. The response is
// accepted/rejected immediately; the outcome is later read from state.
type ActionRequest struct {
	Action     string   `json:"action"`
	Scope      string   `json:"scope,omitempty"`
	Ra         *float64 `json:"ra,omitempty"`
	Dec        *float64 `json:"dec,omitempty"`
	Speed      *int     `json:"speed,omitempty"`
	Direction  string   `json:"direction,omitempty"`
	DurationMs *uint32  `json:"duration_ms,omitempty"`
	Axis       string   `json:"axis,omitempty"`    // move
	Degrees    *float64 `json:"degrees,omitempty"` // move
	Param      string   `json:"param,omitempty"`   // limits
}

// Action sends an operation and returns an error if the mount rejects it.
func (c *Client) Action(req ActionRequest) error {
	return c.request(context.Background(), "action", req, nil)
}
