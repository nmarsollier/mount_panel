//go:build integration

// GOT-limits — behavior at and beyond the axis limits (via goto rejection).
package alpaca

import (
	"net/url"
	"testing"
)

func TestGotoBeyondDecLimitRejected(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	gotoArrival(t, c, 6.0, -60.0)

	// DEC = +80° → dec_axis = ±170°, both outside ±150°.
	r, err := c.put("slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"80.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("goto to DEC=80° was accepted")
	}
	if mustBool(t, c, "slewing") {
		t.Error("rejected goto left slewing active")
	}
}

func TestGotoReturnsToValidAfterRejection(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	gotoArrival(t, c, 6.0, -60.0)

	r, err := c.put("slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"80.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Fatal("out-of-limit goto should fail")
	}

	gotoArrival(t, c, 6.0, -55.0)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false")
	}
}

func TestGotoToExactDecLimit(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	lst := mustFloat(t, c, "siderealtime")
	// RA = LST → HA = 0; DEC=+60 → dec_axis = ±150 (the limit)
	final := gotoArrival(t, c, lst, 60.0)
	assertAngularClose(t, final.ra, final.dec, lst, 60.0, 0.1, "DEC=+60")
}
