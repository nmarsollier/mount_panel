//go:build integration

// STP — Stop and cancellation.
package alpaca

import (
	"net/url"
	"testing"
	"time"
)

func setTime(t *testing.T, c *Client) {
	t.Helper()
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)
}

func TestAbortSlew(t *testing.T) {
	c := newTestClient(t)
	setTime(t, c)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)
	c.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, c)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false after abortslew")
	}
}

func TestStopCancelsQueuedGoto(t *testing.T) {
	c := newTestClient(t)
	setTime(t, c)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"7.0"}, "Declination": {"-50.0"}}, nil)

	c.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, c)

	time.Sleep(5 * time.Second)
	if mustBool(t, c, "slewing") {
		t.Error("the cancelled goto re-appeared after the stop")
	}
}

func TestRepeatStopWhenQuiet(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "abortslew", nil, nil)
	for i := 0; i < 3; i++ {
		r, err := c.put("abortslew", nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		if r.ErrorNumber != 0 {
			t.Errorf("AbortSlew while quiet returned %d", r.ErrorNumber)
		}
	}
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false")
	}
}

func TestMoveAgainAfterStop(t *testing.T) {
	c := newTestClient(t)
	setTime(t, c)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)
	c.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, c)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)
	c.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, c)
}

func TestGlobalStopRest(t *testing.T) {
	t.Skip("requires firmware /api/stop (port 80); the Alpaca equivalent is AbortSlew")
}

func TestGlobalStopClearsTracking(t *testing.T) {
	t.Skip("requires firmware /api/stop (port 80)")
}

func TestAbortSlewPreservesTracking(t *testing.T) {
	c := newTestClient(t)
	setTime(t, c)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)

	c.putOK(t, "abortslew", nil, nil)
	waitTracking(t, c, true)
	if !mustBool(t, c, "tracking") {
		t.Error("AbortSlew cleared tracking")
	}
}

func TestStopThenPark(t *testing.T) {
	c := newTestClient(t)
	setTime(t, c)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)

	c.putOK(t, "abortslew", nil, nil)
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("atpark not true after stop+park")
	}
	c.putOK(t, "unpark", nil, nil)
}
