//go:build integration

// CFG — Site, time and coordinates.
package alpaca

import (
	"math"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSetSite(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)
	c.putOK(t, "siteelevation", url.Values{"SiteElevation": {"750"}}, nil)

	if math.Abs(mustFloat(t, c, "sitelatitude")-(-32.89)) >= 0.01 {
		t.Error("sitelatitude does not match")
	}
	if math.Abs(mustFloat(t, c, "sitelongitude")-(-68.83)) >= 0.01 {
		t.Error("sitelongitude does not match")
	}
	if mustInt(t, c, "siteelevation") != 750 {
		t.Error("siteelevation does not match")
	}
}

func TestSetUTCDate(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	got := mustString(t, c, "utcdate")
	// The bridge reads the PC clock (its authoritative source of time), so the
	// set value is not echoed back like the firmware does; just require a valid
	// ISO 8601 timestamp.
	if !strings.Contains(got, "T") {
		t.Fatalf("UTCDate is not ISO 8601: %s", got)
	}
}

func TestSiderealTimeCoherent(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	lst1 := mustFloat(t, c, "siderealtime")
	time.Sleep(1100 * time.Millisecond)
	lst2 := mustFloat(t, c, "siderealtime")

	if lst1 < 0 || lst1 >= 24 {
		t.Errorf("LST fuera de rango: %v", lst1)
	}
	if math.Mod(lst2-lst1, 24) <= 0 {
		t.Errorf("LST did not advance: %v -> %v", lst1, lst2)
	}
}

func TestChangeSiteDuringMotion(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)

	r, err := c.put("sitelatitude", url.Values{"SiteLatitude": {"abc"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("invalid latitude accepted")
	}
	if math.Abs(mustFloat(t, c, "sitelatitude")-(-32.89)) >= 0.01 {
		t.Error("previous latitude was altered")
	}

	c.putOK(t, "abortslew", nil, nil)
}
