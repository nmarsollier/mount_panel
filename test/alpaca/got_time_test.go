//go:build integration

// GOT-time — a goto's duration must match the slew ramp model.
//
// Skipped: the ramp model is validated against the physical axis distance
// (/api/status), which the Alpaca bridge does not expose. Arrival is covered
// by the other goto tests.
package alpaca

import "testing"

func TestGotoDurationMatchesRamp(t *testing.T) {
	t.Skip("requires the physical axis position to validate the ramp model; not observable via Alpaca")
}

func TestShortGotoDuration(t *testing.T) {
	t.Skip("requires the physical axis position to validate the ramp model; not observable via Alpaca")
}
