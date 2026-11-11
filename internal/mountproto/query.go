package mountproto

import "context"

// DestinationSideOfPierRequest asks the firmware to predict, without moving,
// which pier side a goto to (RA hours, DEC degrees) would select now.
type DestinationSideOfPierRequest struct {
	Query string  `json:"query"`
	Ra    float64 `json:"ra"`
	Dec   float64 `json:"dec"`
}

// DestinationSideOfPierResponse carries the firmware's prediction. When the
// firmware answers ok=false, request() already returned an error, so Side is
// only meaningful on success.
type DestinationSideOfPierResponse struct {
	Type  string `json:"type"`
	ID    uint32 `json:"id"`
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Side  string `json:"side"`
}

// DestinationSideOfPier asks the firmware for the predicted pier side. The
// firmware is the single source of truth: no astronomy is computed here.
func (c *Client) DestinationSideOfPier(ra, dec float64) (*DestinationSideOfPierResponse, error) {
	out := &DestinationSideOfPierResponse{}
	req := DestinationSideOfPierRequest{Query: "destination_side_of_pier", Ra: ra, Dec: dec}
	if err := c.request(context.Background(), "query", req, out); err != nil {
		return nil, err
	}
	return out, nil
}
