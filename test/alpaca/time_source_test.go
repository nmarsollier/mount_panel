//go:build integration

// TIME-SOURCE — UTCDate reports the mount's clock, not the host's.
//
// Fix #7: SetUTCDate writes to the ESP and UTCDate/Status now read that clock
// back instead of returning the host time. A date clearly in the past proves
// which source is being reported. The clock is restored to the host time in
// cleanup so the test leaves no trace.
package alpaca

import (
	"net/url"
	"testing"
	"time"
)

// testUTC is far enough in the past that no host clock would report it, so a
// read-back near it can only come from the mount.
const testUTC = "2020-01-01T00:00:00Z"

func parseUTC(t *testing.T, s string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("utcdate is not ISO 8601 (%q): %v", s, err)
	}
	return ts
}

// setTestClock sets the mount clock to a past date and schedules a cleanup that
// restores the host time, so the test is self-contained.
func setTestClock(t *testing.T, c *Client) {
	t.Helper()
	c.putOK(t, "utcdate", url.Values{"UTCDate": {testUTC}}, nil)
	t.Cleanup(func() {
		now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		c.put("utcdate", url.Values{"UTCDate": {now}}, nil)
	})
}

// TestUTCDateReadsBackMount sets an explicit clock and asserts UTCDate reports
// that mount clock (close to the set value), not the host clock.
func TestUTCDateReadsBackMount(t *testing.T) {
	c := newTestClient(t)
	setTestClock(t, c)

	base := parseUTC(t, testUTC)
	got := parseUTC(t, mustString(t, c, "utcdate"))
	delta := got.Sub(base)
	if delta < 0 || delta > 5*time.Minute {
		t.Fatalf("utcdate = %s, want the mount clock near %s (delta %s)", got, testUTC, delta)
	}
}

// TestUTCDateAdvances confirms the reported clock is a live mount clock that
// moves forward, not a static echo of the host.
func TestUTCDateAdvances(t *testing.T) {
	c := newTestClient(t)
	setTestClock(t, c)

	t1 := parseUTC(t, mustString(t, c, "utcdate"))
	time.Sleep(1500 * time.Millisecond)
	t2 := parseUTC(t, mustString(t, c, "utcdate"))

	if !t2.After(t1) {
		t.Errorf("utcdate did not advance: %s -> %s", t1, t2)
	}
}
