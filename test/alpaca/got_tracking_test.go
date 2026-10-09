//go:build integration

// GOT-tracking — a goto preserves (or leaves off) the tracking state.
package alpaca

import (
	"net/url"
	"testing"
)

func TestGotoPreservesTracking(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	gotoArrival(t, c, 6.0, -60.0)
	if !mustBool(t, c, "tracking") {
		t.Error("goto did not preserve tracking")
	}
}

func TestGotoWithoutTrackingStaysOff(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	gotoArrival(t, c, 6.0, -60.0)
	if mustBool(t, c, "tracking") {
		t.Error("goto left tracking active")
	}
}

func TestGotoPreservesTrackingRate(t *testing.T) {
	for _, rate := range []int{0, 1, 2} {
		t.Run(uitoa(uint32(rate)), func(t *testing.T) {
			c := newTestClient(t)
			setupSiteTime(t, c)
			c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
			c.putOK(t, "trackingrate", url.Values{"TrackingRate": {uitoa(uint32(rate))}}, nil)
			waitTracking(t, c, true)
			gotoArrival(t, c, 6.0, -60.0)
			if !mustBool(t, c, "tracking") {
				t.Error("goto did not preserve tracking")
			}
			if mustInt(t, c, "trackingrate") != rate {
				t.Errorf("goto changed rate from %d to %d", rate, mustInt(t, c, "trackingrate"))
			}
		})
	}
}
