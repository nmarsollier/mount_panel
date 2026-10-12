package mountproto

import "context"

// ControlRequest sets the desired continuous motion. Omitted fields keep their
// current value; a manual rate of zero stops that axis.
type ControlRequest struct {
	Tracking     string   `json:"tracking,omitempty"`
	ManualRaDps  *float64 `json:"manual_ra_dps,omitempty"`
	ManualDecDps *float64 `json:"manual_dec_dps,omitempty"`
}

// Control sends a motion request and returns an error if the mount rejects it.
func (c *Client) Control(req ControlRequest) error {
	return c.request(context.Background(), "control", req, nil)
}
