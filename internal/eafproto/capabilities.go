package eafproto

import "context"

// CapabilitiesResponse is the answer to a capabilities request: the shared
// identification header (device_type, device_id, name, firmware_version) plus
// accessory-specific capabilities, read once at connection time.
type CapabilitiesResponse struct {
	Type            string        `json:"type"`
	ID              uint32        `json:"id"`
	OK              bool          `json:"ok"`
	BootID          string        `json:"boot_id"`
	Proto           int           `json:"proto"`
	DeviceType      string        `json:"device_type"`
	DeviceID        string        `json:"device_id"`
	Name            string        `json:"name"`
	FirmwareVersion string        `json:"firmware_version"`
	Capabilities    AccessoryCaps `json:"capabilities"`
}

// AccessoryCaps is the accessory-specific part of the capabilities response.
type AccessoryCaps struct {
	Focuser struct {
		RelativeMove  bool `json:"relative_move"`
		SignedCounter bool `json:"signed_counter"`
		Halt          bool `json:"halt"`
		LocalButtons  bool `json:"local_buttons"`
		TMCUART       bool `json:"tmc_uart"`
		Encoder       bool `json:"encoder"`
		HomeSensor    bool `json:"home_sensor"`
		TempComp      bool `json:"temp_comp"`
	} `json:"focuser"`
	Fan struct {
		Modes []string `json:"modes"`
		PWM   bool     `json:"pwm"`
		Relay bool     `json:"relay"`
		Tach  bool     `json:"tach"`
	} `json:"fan"`
	Sensors struct {
		Ambient   string `json:"ambient"`
		Mirror    string `json:"mirror"`
		Vibration string `json:"vibration"`
	} `json:"sensors"`
}

// Capabilities asks the accessory for its identity and static capabilities.
func (c *Client) Capabilities() (*CapabilitiesResponse, error) {
	out := &CapabilitiesResponse{}
	if err := c.request(context.Background(), "capabilities", "", "", nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
