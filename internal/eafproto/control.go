package eafproto

import "context"

// FanSetPower switches the fan's 12 V relay (the power channel).
func (c *Client) FanSetPower(on bool) error {
	return c.request(context.Background(), "control", "fan", "power",
		map[string]any{"value": on}, nil)
}

// FanSetSpeed sets the fan PWM duty (0..100 %).
func (c *Client) FanSetSpeed(pct int) error {
	return c.request(context.Background(), "control", "fan", "speed",
		map[string]any{"value": pct}, nil)
}
