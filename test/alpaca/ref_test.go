//go:build integration

// REF — Reference, home and position recovery.
package alpaca

import (
	"net/url"
	"testing"
)

func TestHomeParkCelestialDistinct(t *testing.T) {
	c := newTestClient(t)
	// freshly started: AtHome/AtPark are booleans reflecting the real state
	if _, err := c.boolValue("athome", nil); err != nil {
		t.Errorf("athome: %v", err)
	}
	if _, err := c.boolValue("atpark", nil); err != nil {
		t.Errorf("atpark: %v", err)
	}
}

func TestFindHomeFromPosition(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)

	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false after findhome")
	}
	if !mustBool(t, c, "athome") {
		t.Error("athome not true after findhome")
	}
}

func TestFindHomeWhenAlreadyHome(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)
	if !mustBool(t, c, "athome") {
		t.Error("athome not true")
	}
}

func TestFindHomeFromVariousPositions(t *testing.T) {
	for _, pos := range [][2]float64{{6.0, -60.0}, {8.0, -45.0}, {10.0, -70.0}} {
		t.Run(fs(pos[0])+"h/"+fs(pos[1]), func(t *testing.T) {
			c := newTestClient(t)
			c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
			c.putOK(t, "slewtocoordinatesasync",
				url.Values{"RightAscension": {fs(pos[0])}, "Declination": {fs(pos[1])}}, nil)
			waitNotSlewing(t, c)

			c.putOK(t, "findhome", nil, nil)
			waitNotSlewing(t, c)
			if !mustBool(t, c, "athome") {
				t.Error("athome not true (axis position not verifiable via Alpaca)")
			}
		})
	}
}

func TestFindHomeInTracking(t *testing.T) {
	for _, rate := range []int{0, 1, 2} {
		t.Run(uitoa(uint32(rate)), func(t *testing.T) {
			c := newTestClient(t)
			c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
			c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
			c.putOK(t, "trackingrate", url.Values{"TrackingRate": {uitoa(uint32(rate))}}, nil)
			waitTracking(t, c, true)

			c.putOK(t, "findhome", nil, nil)
			waitNotSlewing(t, c)
			if !mustBool(t, c, "athome") {
				t.Error("athome not true")
			}
			if mustBool(t, c, "tracking") {
				t.Error("home should not resume tracking")
			}
		})
	}
}

func TestFindHomeMidGoto(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)

	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)
	if !mustBool(t, c, "athome") {
		t.Error("home mid-goto did not return home")
	}
}
