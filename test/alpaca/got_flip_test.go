//go:build integration

// GOT-flip — a goto across the meridian flips the pier side.
package alpaca

import (
	"testing"
)

func TestMeridianFlip(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")

	// HA = +4h (west) → pierEast
	raWest := mathMod(lst-4.0, 24.0)
	gotoArrival(t, c, raWest, -60.0)
	if mustInt(t, c, "sideofpier") != 0 {
		t.Errorf("HA=+4h should be pierEast, got %d", mustInt(t, c, "sideofpier"))
	}

	// HA = -4h (east) → pierWest (the flip)
	raEast := mathMod(lst+4.0, 24.0)
	gotoArrival(t, c, raEast, -60.0)
	if mustInt(t, c, "sideofpier") != 1 {
		t.Errorf("HA=-4h should be pierWest (flip), got %d", mustInt(t, c, "sideofpier"))
	}
}

func TestGotoNearMeridian(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")
	ra := mathMod(lst, 24.0) // HA ≈ 0
	dec := -60.0

	final := gotoArrival(t, c, ra, dec)
	assertAngularClose(t, final.ra, final.dec, ra, dec, reportedPosDeg, "near-meridian goto")
}
