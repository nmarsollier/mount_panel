//go:build integration

// GOT-sites — the same goto resolves correctly from different site coordinates.
package alpaca

import (
	"net/url"
	"testing"
)

func TestGotoAcrossSites(t *testing.T) {
	sites := [][2]string{
		{"-32.89", "-68.83"}, // southern, home site
		{"-23.55", "-46.63"}, // southern, different longitude
		{"40.71", "-74.01"},  // northern
	}
	for _, site := range sites {
		t.Run(site[0]+"/"+site[1], func(t *testing.T) {
			c := newTestClient(t)
			c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
			c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {site[0]}}, nil)
			c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {site[1]}}, nil)
			final := gotoArrival(t, c, 6.0, -60.0)
			assertAngularClose(t, final.ra, final.dec, 6.0, -60.0, 0.1, "site "+site[0]+"/"+site[1])
		})
	}
}
