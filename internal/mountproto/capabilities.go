package mountproto

import "context"

// CapabilitiesResponse is the answer to a capabilities request: the shared
// identification header (device_type, device_id, name, firmware_version) plus
// mount-specific capabilities, read once at connection time.
type CapabilitiesResponse struct {
	Type            string    `json:"type"`
	ID              uint32    `json:"id"`
	OK              bool      `json:"ok"`
	Error           string    `json:"error"`
	BootID          string    `json:"boot_id"`
	Proto           int       `json:"proto"`
	DeviceType      string    `json:"device_type"`
	DeviceID        string    `json:"device_id"`
	Name            string    `json:"name"`
	FirmwareVersion string    `json:"firmware_version"`
	Capabilities    MountCaps `json:"capabilities"`
}

// MountCaps is the mount-specific part of the capabilities response.
type MountCaps struct {
	Axes              []string  `json:"axes"`
	SlewSpeedsDps     []float64 `json:"slew_speeds_dps"`
	MoveAxisSpeedsDps []float64 `json:"move_axis_speeds_dps"`
	TrackingModes     []string  `json:"tracking_modes"`
	Caps              struct {
		Guide    bool `json:"guide"`
		Home     bool `json:"home"`
		Park     bool `json:"park"`
		Unpark   bool `json:"unpark"`
		MoveAxis bool `json:"move_axis"`
		Slew     bool `json:"slew"`
		Tracking bool `json:"tracking"`
		Sync     bool `json:"sync"`
		// MeridianFlip is a pointer so an absent field (old firmware) is
		// distinguishable from an explicit false (fork).
		MeridianFlip *bool `json:"meridian_flip"`
		PierSide     bool  `json:"pier_side"`
	} `json:"caps"`
}

// Capabilities asks the mount for its identity and static capabilities.
func (c *Client) Capabilities() (*CapabilitiesResponse, error) {
	out := &CapabilitiesResponse{}
	if err := c.request(context.Background(), "capabilities", nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
