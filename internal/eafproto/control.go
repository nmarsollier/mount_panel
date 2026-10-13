package eafproto

import "context"

// FanSetMode sets the fan mode: "off", "auto" or "manual".
func (c *Client) FanSetMode(mode string) error {
	return c.request(context.Background(), "control", "fan", "mode",
		map[string]any{"value": mode}, nil)
}

// FanSetSpeed sets the manual fan speed (1..100). Applied immediately in
// MANUAL mode and remembered otherwise.
func (c *Client) FanSetSpeed(speed int) error {
	return c.request(context.Background(), "control", "fan", "speed",
		map[string]any{"value": speed}, nil)
}

// FocuserSetEnabled enables or disables the focuser driver. Disabling releases
// the TMC2209 EN pin (no torque) and rejects motion until re-enabled.
func (c *Client) FocuserSetEnabled(enabled bool) error {
	return c.request(context.Background(), "control", "focuser", "enabled",
		map[string]any{"value": enabled}, nil)
}
