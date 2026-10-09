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
	// just west (HA=+1h) → pierEast
	raW := mathMod(lst-1.0, 24.0)
	gotoArrival(t, c, raW, -60.0)
	if mustInt(t, c, "sideofpier") != 0 {
		t.Error("expected pierEast")
	}
	// correct to just east (HA=-1h) → pierWest (flip)
	raE := mathMod(lst+1.0, 24.0)
	gotoArrival(t, c, raE, -60.0)
	if mustInt(t, c, "sideofpier") != 1 {
		t.Error("expected pierWest (flip)")
	}
}
