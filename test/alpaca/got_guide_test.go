//go:build integration

// GOT-guide — a PulseGuide issued mid-goto must not corrupt the slew.
package alpaca

import (
	"net/url"
	"testing"
)

func TestPulseGuideDuringGoto(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)

	// home first so the following goto is always a large, in-flight move
	c.putOK(t, "findhome", nil, nil)
	waitNotSlewing(t, c)

	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)

	// pulse guide while the slew is in flight (may be deferred/rejected)
	_, _ = c.put("pulseguide", url.Values{"Direction": {"0"}, "Duration": {"200"}}, nil)

	final := waitNotSlewing(t, c)
	if final.slewing {
		t.Error("the goto did not finish after the pulse guide")
	}
	if mustBool(t, c, "ispulseguiding") {
		t.Error("left pulse-guiding active")
	}
}
