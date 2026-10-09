//go:build integration

// MER — Pier side and meridian.
package alpaca

import (
	"net/url"
	"testing"
)

func TestSideOfPierReport(t *testing.T) {
	c := newTestClient(t)
	sop := mustInt(t, c, "sideofpier")
	if sop != 0 && sop != 1 {
		t.Errorf("invalid SideOfPier: %d", sop)
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
	if v := asInt(r.Value); v != 0 && v != 1 {
		t.Errorf("invalid DestinationSideOfPier: %v", r.Value)
	}
}

func TestMeridianFlipSameObject(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)

	lst := mustFloat(t, c, "siderealtime")
	ra := mathMod(lst+2.0, 24.0) // HA = -2h (east) → pierWest

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)
	side1 := mustInt(t, c, "sideofpier")
	if side1 != 1 {
		t.Errorf("HA=-2h should be pierWest, got %d", side1)
	}

	// advance the clock 4h so this RA crosses the meridian (HA → +2h)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T06:00:00"}}, nil)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)
	side2 := mustInt(t, c, "sideofpier")
	if side2 != 0 {
		t.Errorf("HA=+2h should be pierEast (flip), got %d", side2)
	}
}

func TestMeridianFlipPreservesTracking(t *testing.T) {
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
		t.Error("flip GOTO did not resume tracking")
	}
}
