//go:build integration

// GOT — Pointing at real objects (slew / goto).
package alpaca

import (
	"net/url"
	"testing"
)

func TestGotoFromRest(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	ra, dec := 6.0, -60.0
	final := gotoArrival(t, c, ra, dec)
	assertAngularClose(t, final.ra, final.dec, ra, dec, reportedPosDeg, "slew from rest")
}

func TestTargetAndSlewToTarget(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	ra, dec := 8.0, -45.0
	c.putOK(t, "targetrightascension", url.Values{"TargetRightAscension": {fs(ra)}}, nil)
	c.putOK(t, "targetdeclination", url.Values{"TargetDeclination": {fs(dec)}}, nil)
	if mustFloat(t, c, "targetrightascension") != ra {
		t.Error("targetrightascension does not match")
	}
	if mustFloat(t, c, "targetdeclination") != dec {
		t.Error("targetdeclination does not match")
	}

	c.putOK(t, "slewtotargetasync", nil, nil)
	final := waitNotSlewing(t, c)
	assertAngularClose(t, final.ra, final.dec, ra, dec, reportedPosDeg, "apuntado a objetivo")
}

func TestGotoToCurrentPosition(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	ra := mustFloat(t, c, "rightascension")
	dec := mustFloat(t, c, "declination")
	gotoArrival(t, c, ra, dec)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false")
	}
}

func TestGotoImpossible(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)

	r, err := c.put("slewtocoordinatesasync",
		url.Values{"RightAscension": {"25.0"}, "Declination": {"0.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("GOTO with RA=25 accepted")
	}
	if mustBool(t, c, "slewing") {
		t.Error("the rejected GOTO started slewing")
	}
}

func TestGotoDuringGoto(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"8.0"}, "Declination": {"-45.0"}}, nil)
	final := waitNotSlewing(t, c)
	if final.slewing {
		t.Error("final slewing not false")
	}
}
