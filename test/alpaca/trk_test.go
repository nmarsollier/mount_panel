//go:build integration

// TRK — Tracking during observation and exposure.
package alpaca

import (
	"math"
	"net/url"
	"testing"
	"time"
)

func raWrapDiff(a, b float64) float64 {
	d := math.Abs(a - b)
	return math.Min(d, 24.0-d)
}

func TestDisableTracking(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
	waitTracking(t, c, false)
	if mustBool(t, c, "tracking") {
		t.Error("tracking not false")
	}
}

func TestReactivateTracking(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
	waitTracking(t, c, false)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	if !mustBool(t, c, "tracking") {
		t.Error("tracking not true")
	}
}

func TestChangeTrackingRate(t *testing.T) {
	for _, rate := range []int{0, 1, 2} {
		t.Run(uitoa(uint32(rate)), func(t *testing.T) {
			c := newTestClient(t)
			c.putOK(t, "trackingrate", url.Values{"TrackingRate": {uitoa(uint32(rate))}}, nil)
			if mustInt(t, c, "trackingrate") != rate {
				t.Error("trackingrate does not match")
			}
		})
	}
}

func TestTrackingOrderIdempotent(t *testing.T) {
	c := newTestClient(t)
	for i := 0; i < 5; i++ {
		c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
		waitTracking(t, c, true)
	}
	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
	waitTracking(t, c, false)
}

func TestTrackingKeepsRAStable(t *testing.T) {
	c := newTestClient(t)

	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
	waitTracking(t, c, false)
	ra0 := mustFloat(t, c, "rightascension")
	time.Sleep(3 * time.Second)
	ra1 := mustFloat(t, c, "rightascension")
	driftOff := raWrapDiff(ra0, ra1)

	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	ra2 := mustFloat(t, c, "rightascension")
	time.Sleep(3 * time.Second)
	ra3 := mustFloat(t, c, "rightascension")
	driftOn := raWrapDiff(ra2, ra3)

	if driftOff <= 0.0003 {
		t.Errorf("baseline drift too small to measure: %f h", driftOff)
	}
	if driftOn >= driftOff*0.25 {
		t.Errorf("tracking no mantuvo RA: off=%f h on=%f h", driftOff, driftOn)
	}
}

func TestTrackingProlongedWithQueries(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)
	for i := 0; i < 10; i++ {
		if !mustBool(t, c, "tracking") {
			t.Fatal("tracking not true")
		}
		if mustBool(t, c, "slewing") {
			t.Fatal("unexpected slewing during tracking")
		}
		time.Sleep(500 * time.Millisecond)
	}
	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
}

func TestTrackingRateCoherence(t *testing.T) {
	t.Skip("requires reading the physical axis position (ra_deg); not observable via Alpaca")
}

func TestTrackingStopsAtLimit(t *testing.T) {
	t.Skip("requires reading the physical axis position and limits; not observable via Alpaca")
}
