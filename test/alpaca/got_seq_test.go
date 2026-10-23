//go:build integration

// GOT-sequence — consecutive gotos and NINA-style correction.
package alpaca

import (
	"testing"
)

func TestConsecutiveGotos(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	for _, pos := range [][2]float64{{6.0, -60.0}, {8.0, -45.0}, {10.0, -70.0}, {6.0, -55.0}} {
		final := gotoArrival(t, c, pos[0], pos[1])
		assertAngularClose(t, final.ra, final.dec, pos[0], pos[1], 0.1, "goto")
	}
}

func TestNinaStyleGotoThenCorrect(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	// first a long goto, then a small correction
	gotoArrival(t, c, 6.0, -60.0)
	gotoArrival(t, c, 6.0, -62.0)
}

func TestGotoFromHome(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)
	if !mustBool(t, c, "athome") {
		t.Fatal("home is not (0,0)")
	}
	final := gotoArrival(t, c, 6.0, -60.0)
	assertAngularClose(t, final.ra, final.dec, 6.0, -60.0, 0.1, "from home")
}

func TestNinaCorrectAtDecLimit(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")
	// near the limit (DEC=+58), then a small correction
	gotoArrival(t, c, lst, 58.0)
	gotoArrival(t, c, lst, 55.0)
}

func TestNinaCorrectAcrossMeridian(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")
	gem := germanEquatorial(t, c)
	// just west (HA=+1h)
	raW := mathMod(lst-1.0, 24.0)
	gotoArrival(t, c, raW, -60.0)
	sideWest := mustInt(t, c, "sideofpier")
	// correct to just east (HA=-1h)
	raE := mathMod(lst+1.0, 24.0)
	gotoArrival(t, c, raE, -60.0)
	sideEast := mustInt(t, c, "sideofpier")

	if gem {
		if sideWest != pierEast {
			t.Errorf("HA=+1h side = %d, want pierEast(0)", sideWest)
		}
		if sideEast != pierWest {
			t.Errorf("HA=-1h side = %d, want pierWest(1) (flip)", sideEast)
		}
	} else {
		if sideWest != pierUnknown {
			t.Errorf("HA=+1h side = %d, want pierUnknown(-1)", sideWest)
		}
		if sideEast != pierUnknown {
			t.Errorf("HA=-1h side = %d, want pierUnknown(-1)", sideEast)
		}
	}
}
