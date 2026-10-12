package eafproto

import "context"

// FanModeResponse is the fan state after a mode change.
type FanModeResponse struct {
	Mode         string `json:"mode"`
	ModeRev      int    `json:"mode_rev"`
	OutputPWMPct int    `json:"output_pwm_pct"`
	RelayCommand bool   `json:"relay_command"`
	Reason       string `json:"reason"`
}

// FanSetMode selects the fan mode: "normal", "soft" or "off".
func (c *Client) FanSetMode(mode string) (*FanModeResponse, error) {
	out := struct {
		Data FanModeResponse `json:"data"`
	}{}
	err := c.request(context.Background(), "control", "fan", "mode",
		map[string]any{"mode": mode}, &out)
	if err != nil {
		return nil, err
	}
	return &out.Data, nil
}
