package proto

import "context"

// CapabilitiesResponse is the answer to a capabilities request: identity and
// static capabilities, read once at connection time.
type CapabilitiesResponse struct {
	Type          string    `json:"type"`
	ID            uint32    `json:"id"`
	OK            bool      `json:"ok"`
	Error         string    `json:"error"`
	BootID        uint32    `json:"boot_id"`
	Proto         int       `json:"proto"`
	Name          string    `json:"name"`
	Axes          []string  `json:"axes"`
	SlewSpeedsDps []float64 `json:"slew_speeds_dps"`
	TrackingModes []string  `json:"tracking_modes"`
	Caps          struct {
		Guide    bool `json:"guide"`
		Home     bool `json:"home"`
		Park     bool `json:"park"`
		Unpark   bool `json:"unpark"`
		MoveAxis bool `json:"move_axis"`
		Slew     bool `json:"slew"`
		Tracking bool `json:"tracking"`
		Sync     bool `json:"sync"`
	} `json:"caps"`
	ConfigRev uint32 `json:"config_rev"`
}

// Capabilities asks the mount for its identity and static capabilities.
func (c *Client) Capabilities() (*CapabilitiesResponse, error) {
	out := &CapabilitiesResponse{}
	if err := c.request(context.Background(), "capabilities", nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
