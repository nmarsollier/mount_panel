package mountproto

import "context"

// StateResponse is a coherent snapshot of the whole mount, as reported by the
// firmware in a single read.
type StateResponse struct {
	Type        string  `json:"type"`
	ID          uint32  `json:"id"`
	OK          bool    `json:"ok"`
	Error       string  `json:"error"`
	State       string  `json:"state"`
	Tracking    string  `json:"tracking"`
	RaSteps     int64   `json:"ra_steps"`
	DecSteps    int64   `json:"dec_steps"`
	RaAxisDeg   float64 `json:"ra_axis_deg"`
	DecAxisDeg  float64 `json:"dec_axis_deg"`
	Ra          float64 `json:"ra"`
	Dec         float64 `json:"dec"`
	Lst         float64 `json:"lst"`
	PierSide    string  `json:"pier_side"`
	RaSpeedDps  float64 `json:"ra_speed_dps"`
	DecSpeedDps float64 `json:"dec_speed_dps"`
	Guiding     bool    `json:"guiding"`
	AtHome      bool    `json:"at_home"`
	AtPark      bool    `json:"at_park"`
	Power       bool    `json:"power"`
	TimeValid   bool    `json:"time_valid"`
	Limits      struct {
		RaMin  float64 `json:"ra_min"`
		RaMax  float64 `json:"ra_max"`
		DecMin float64 `json:"dec_min"`
		DecMax float64 `json:"dec_max"`
	} `json:"limits"`
	ConfigRev uint32 `json:"config_rev"`
	UptimeMs  uint64 `json:"uptime_ms"`
}

// State asks the mount for a coherent snapshot.
func (c *Client) State() (*StateResponse, error) {
	out := &StateResponse{}
	if err := c.request(context.Background(), "state", nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
