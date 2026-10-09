//go:build integration

// GUIDE-effect — a PulseGuide moves the axis by the expected amount and time.
//
// Skipped: these validate the physical axis step change (rate × duration)
// against the raw step counters and the X-NRM-Snapshot header, both of which
// are firmware test-mode observability not exposed by the Alpaca bridge. The
// PulseGuide acceptance contract is covered in gui_test.go.
package alpaca

import "testing"

func TestPulseEffectMagnitudeAndTime(t *testing.T) {
	t.Skip("requires physical axis position (steps) and snapshot header; not observable via Alpaca")
}

func TestPulseEastMovesRANegative(t *testing.T) {
	t.Skip("requires physical axis position; not observable via Alpaca")
}

func TestPulseFlipReversesDec(t *testing.T) {
	t.Skip("requires physical axis position; not observable via Alpaca")
}

func TestPulseTrackingModes(t *testing.T) {
	t.Skip("requires physical axis position; not observable via Alpaca")
}

func TestPulseAtDecLimitClamped(t *testing.T) {
	t.Skip("requires physical axis position and limits; not observable via Alpaca")
}
