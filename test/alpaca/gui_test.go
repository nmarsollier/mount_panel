//go:build integration

// GUI — Autoguiding and fine corrections (PulseGuide).
package alpaca

import (
	"net/url"
	"testing"
	"time"
)

func TestPulseGuideDirections(t *testing.T) {
	for _, dir := range []int{0, 1, 2, 3} {
		t.Run(string(rune('0'+dir)), func(t *testing.T) {
			c := newTestClient(t)
			c.putOK(t, "pulseguide", url.Values{"Direction": {uitoa(uint32(dir))}, "Duration": {"100"}}, nil)
			if v, err := c.boolValue("ispulseguiding", nil); err != nil {
				t.Fatal(err)
			} else if v != true && v != false {
				t.Errorf("invalid ispulseguiding: %v", v)
			}
		})
	}
}

func TestPulseDurationControlsCorrection(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "pulseguide", url.Values{"Direction": {"0"}, "Duration": {"300"}}, nil)
	time.Sleep(500 * time.Millisecond)
	if mustBool(t, c, "ispulseguiding") {
		t.Error("ispulseguiding still true after the pulse")
	}
}

func TestSetGuideRates(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "guideraterightascension", url.Values{"GuideRateRightAscension": {"0.001"}}, nil)
	c.putOK(t, "guideratedeclination", url.Values{"GuideRateDeclination": {"0.001"}}, nil)
	if mustFloat(t, c, "guideraterightascension") <= 0 {
		t.Error("guideraterightascension no > 0")
	}
	if mustFloat(t, c, "guideratedeclination") <= 0 {
		t.Error("guideratedeclination no > 0")
	}
}

func TestGuideRAAndDec(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "pulseguide", url.Values{"Direction": {"0"}, "Duration": {"200"}}, nil)
	c.putOK(t, "pulseguide", url.Values{"Direction": {"2"}, "Duration": {"200"}}, nil)
	// no error when overlapping distinct axes
}

func TestGuideWhenParked(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "park", nil, nil)
	r, err := c.put("pulseguide", url.Values{"Direction": {"0"}, "Duration": {"100"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("PulseGuide while parked was accepted")
	}
	c.putOK(t, "unpark", nil, nil)
}
