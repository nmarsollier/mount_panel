package mount

import "testing"

// TestMovementIntervals checks that consecutive moving samples compress into one
// interval and that a stop breaks the run.
func TestMovementIntervals(t *testing.T) {
	samples := []MovementSample{
		{TsMs: 1000, Moving: true},
		{TsMs: 2000, Moving: true},
		{TsMs: 3000, Moving: false},
		{TsMs: 4000, Moving: true},
		{TsMs: 5000, Moving: true},
	}
	got := movementIntervals(samples)

	if len(got) != 2 {
		t.Fatalf("got %d intervals, want 2: %+v", len(got), got)
	}
	if got[0].StartMs != 1000 || got[0].EndMs != 2000 {
		t.Fatalf("first interval = %+v, want [1000, 2000]", got[0])
	}
	if got[1].StartMs != 4000 || got[1].EndMs != 5000 {
		t.Fatalf("second interval = %+v, want [4000, 5000]", got[1])
	}
}

// TestMovementIntervalsBridgesHole checks that a missed poll (no sample between
// two moving ones) is bridged into the same run.
func TestMovementIntervalsBridgesHole(t *testing.T) {
	samples := []MovementSample{
		{TsMs: 1000, Moving: true},
		{TsMs: 3000, Moving: true},
	}
	got := movementIntervals(samples)
	if len(got) != 1 || got[0].StartMs != 1000 || got[0].EndMs != 3000 {
		t.Fatalf("hole not bridged: %+v", got)
	}
}
