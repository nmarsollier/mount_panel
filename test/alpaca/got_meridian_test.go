//go:build integration

// GOT-meridian — a goto across the meridian. A GEM flips (pierEast ↔ pierWest);
// a fork keeps its orientation (pierUnknown).
package alpaca

import (
	"testing"
)

func TestMeridianCrossingGoto(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")
	gem := germanEquatorial(t, c)

	// HA = +4h (west)
	raWest := mathMod(lst-4.0, 24.0)
	gotoArrival(t, c, raWest, -60.0)
	sideWest := mustInt(t, c, "sideofpier")

	// HA = -4h (east)
	raEast := mathMod(lst+4.0, 24.0)
	gotoArrival(t, c, raEast, -60.0)
	sideEast := mustInt(t, c, "sideofpier")

	if gem {
		if sideWest != pierEast {
			t.Errorf("HA=+4h side = %d, want pierEast(0)", sideWest)
		}
		if sideEast != pierWest {
			t.Errorf("HA=-4h side = %d, want pierWest(1) (flip)", sideEast)
		}
	} else {
		if sideWest != pierUnknown {
			t.Errorf("HA=+4h side = %d, want pierUnknown(-1)", sideWest)
		}
		if sideEast != pierUnknown {
			t.Errorf("HA=-4h side = %d, want pierUnknown(-1)", sideEast)
		}
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
