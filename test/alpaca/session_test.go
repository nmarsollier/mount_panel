//go:build integration

// SESSION — long mixed sessions and interior-limit behaviour.
package alpaca

import (
	"net/url"
	"testing"
)

func TestTrackingStopsAtTestLimit(t *testing.T) {
	t.Skip("requires changing limits at runtime (firmware /api/limits), not exposed by the bridge")
}

func TestProlongedMixedSession(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	other := newClient()

	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	waitTracking(t, c, true)

	// a guide pulse while tracking
	c.putOK(t, "guideratedeclination", url.Values{"GuideRateDeclination": {"0.5"}}, nil)
	c.putOK(t, "pulseguide", url.Values{"Direction": {"0"}, "Duration": {"300"}}, nil)

	// a small goto preserves tracking
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitNotSlewing(t, c)
	if !mustBool(t, c, "tracking") {
		t.Error("tracking lost after goto")
	}

	// a second client observes without disturbing anything
	if mustBool(t, other, "slewing") {
		t.Error("second client sees unexpected slewing")
	}
	if !mustBool(t, other, "tracking") {
		t.Error("second client does not see tracking")
	}

	// cancel and park cleanly
	c.putOK(t, "abortslew", nil, nil)
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("atpark not true after park")
	}
	c.putOK(t, "unpark", nil, nil)
	if mustBool(t, c, "atpark") {
		t.Error("atpark not false after unpark")
	}

	// return to a known state so later tests start from home
	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)
}
