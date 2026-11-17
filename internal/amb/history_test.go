package amb

import "testing"

// TestVibrationStats checks the summary math (mean/min/max/stddev/last) over a
// known window.
func TestVibrationStats(t *testing.T) {
	samples := []VibrationSample{
		{TsMs: 1000, VelRms: 2},
		{TsMs: 2000, VelRms: 4},
		{TsMs: 3000, VelRms: 6},
	}
	got := vibrationStats(samples)

	if got.Count != 3 {
		t.Fatalf("count = %d, want 3", got.Count)
	}
	if *got.Min != 2 || *got.Max != 6 || *got.Mean != 4 || *got.Last != 6 {
		t.Fatalf("unexpected stats: %+v", got)
	}
	// stddev = sqrt(((2-4)^2 + (4-4)^2 + (6-4)^2) / 3) = sqrt(8/3) ≈ 1.633
	if *got.Stddev < 1.63 || *got.Stddev > 1.64 {
		t.Fatalf("stddev = %v, want ~1.633", *got.Stddev)
	}
}

// TestVibrationStatsEmpty checks that an empty window yields nil stat fields.
func TestVibrationStatsEmpty(t *testing.T) {
	got := vibrationStats(nil)
	if got.Count != 0 || got.Min != nil || got.Max != nil || got.Mean != nil || got.Stddev != nil || got.Last != nil {
		t.Fatalf("expected zero stats, got %+v", got)
	}
}
