//go:build integration

// MER — Pier side and meridian. A German equatorial (alignment mode 2) flips
// across the meridian and reports pierEast/pierWest; a fork (mode 1) never
// flips and always reports pierUnknown.
package alpaca

import (
	"net/url"
	"testing"
)

func TestSideOfPierReport(t *testing.T) {
	c := newTestClient(t)
	sop := mustInt(t, c, "sideofpier")
	if germanEquatorial(t, c) {
		if sop != pierEast && sop != pierWest {
			t.Errorf("SideOfPier = %d, want pierEast(0) or pierWest(1)", sop)
		}
	} else if sop != pierUnknown {
		t.Errorf("SideOfPier = %d, want pierUnknown(-1)", sop)
	}
}

func TestDestinationSideOfPier(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	r, err := c.get("destinationsideofpier",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber != 0 {
		t.Fatalf("DestinationSideOfPier err=%d", r.ErrorNumber)
	}
	v := asInt(r.Value)
	if germanEquatorial(t, c) {
		if v != pierEast && v != pierWest {
			t.Errorf("DestinationSideOfPier = %v, want pierEast(0) or pierWest(1)", r.Value)
		}
	} else if v != pierUnknown {
		t.Errorf("DestinationSideOfPier = %v, want pierUnknown(-1)", r.Value)
	}
}

// TestMeridianCrossing checks the pier side before and after a target crosses
// the meridian. A GEM flips (pierWest → pierEast); a fork reports pierUnknown.
func TestMeridianCrossing(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)

	gem := germanEquatorial(t, c)

	lst := mustFloat(t, c, "siderealtime")
	ra := mathMod(lst+2.0, 24.0) // HA = -2h (east)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)
	side1 := mustInt(t, c, "sideofpier")

	// advance the clock 4h so this RA crosses the meridian (HA → +2h)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T06:00:00"}}, nil)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)
	side2 := mustInt(t, c, "sideofpier")

	if gem {
		if side1 != pierWest {
			t.Errorf("HA=-2h side = %d, want pierWest(1)", side1)
		}
		if side2 != pierEast {
			t.Errorf("HA=+2h side = %d, want pierEast(0) (flip)", side2)
		}
	} else {
		if side1 != pierUnknown {
			t.Errorf("side = %d, want pierUnknown(-1)", side1)
		}
		if side2 != pierUnknown {
			t.Errorf("side after crossing = %d, want pierUnknown(-1)", side2)
		}
	}
}

func TestMeridianCrossingPreservesTracking(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)

	lst := mustFloat(t, c, "siderealtime")
	ra := mathMod(lst+2.0, 24.0)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)

	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)

	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T06:00:00"}}, nil)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)

	if !mustBool(t, c, "tracking") {
		t.Error("GOTO did not resume tracking")
	}
}
