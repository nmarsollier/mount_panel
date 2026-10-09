package proto

import "context"

// ConfigResponse is the current configuration echoed by the mount. A request
// with no fields is a read; present fields are updated and echoed back.
type ConfigResponse struct {
	Type      string `json:"type"`
	ID        uint32 `json:"id"`
	OK        bool   `json:"ok"`
	Error     string `json:"error"`
	ConfigRev uint32 `json:"config_rev"`
	Config    struct {
		Lat          float64 `json:"lat"`
		Lon          float64 `json:"lon"`
		Elevation    int     `json:"elevation"`
		GuideRateRa  float64 `json:"guide_rate_ra"`
		GuideRateDec float64 `json:"guide_rate_dec"`
	} `json:"config"`
}

// ConfigRequest carries the optional fields of a config write; a nil pointer
// leaves that field untouched.
type ConfigRequest struct {
	Lat          *float64 `json:"lat,omitempty"`
	Lon          *float64 `json:"lon,omitempty"`
	Elevation    *int     `json:"elevation,omitempty"`
	GuideRateRa  *float64 `json:"guide_rate_ra,omitempty"`
	GuideRateDec *float64 `json:"guide_rate_dec,omitempty"`
	Utc          *string  `json:"utc,omitempty"`
}

// ConfigGet reads the current configuration.
func (c *Client) ConfigGet() (*ConfigResponse, error) {
	out := &ConfigResponse{}
	if err := c.request(context.Background(), "config", nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigSet writes the given configuration fields and returns the updated one.
func (c *Client) ConfigSet(req ConfigRequest) (*ConfigResponse, error) {
	out := &ConfigResponse{}
	if err := c.request(context.Background(), "config", req, out); err != nil {
		return nil, err
	}
	return out, nil
}
