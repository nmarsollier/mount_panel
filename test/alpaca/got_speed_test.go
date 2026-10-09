//go:build integration

// GOT-speed — goto duration scales with the requested speed (firmware ramp).
//
// These cases validate the firmware's motion ramp model from the physical axis
// distance, which requires the debug REST endpoint /api/slew-to-coordinates
// (to select the speed profile) and /api/status (raw axis position). Neither
// is part of the Alpaca bridge, so the whole file is skipped.
package alpaca

import "testing"

func TestGotoSpeedDuration(t *testing.T) {
	t.Skip("requires /api/slew-to-coordinates (speed selection) and physical axis position; not observable via Alpaca")
}

func TestSmallGotoTime(t *testing.T) {
	t.Skip("requires the physical axis position to validate the ramp model; not observable via Alpaca")
}
