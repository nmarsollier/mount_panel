//go:build integration

// VAL — Public contract, invalid inputs and observable errors.
package alpaca

import (
	"net/url"
	"strconv"
	"testing"
)

func TestRejectInvalidNumericOrEnum(t *testing.T) {
	cases := []struct{ name, key, value string }{
		{"sitelatitude", "SiteLatitude", "NaN"},
		{"sitelongitude", "SiteLongitude", "Infinity"},
		{"siteelevation", "SiteElevation", "abc"},
		{"utcdate", "UTCDate", "not-a-date"},
		{"tracking", "Tracking", "bogus"},
		{"trackingrate", "TrackingRate", "99"},
		{"declinationrate", "DeclinationRate", "NaN"},
		{"rightascensionrate", "RightAscensionRate", "-Infinity"},
		{"slewsettletime", "SlewSettleTime", "-5"},
		{"targetrightascension", "TargetRightAscension", "30.0"},
		{"targetdeclination", "TargetDeclination", "100.0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestClient(t)
			r, err := c.put(tc.name, url.Values{tc.key: {tc.value}}, nil)
			if err != nil {
				t.Fatalf("%s %v: %v", tc.name, tc.value, err)
			}
			if r.ErrorNumber == 0 {
				t.Errorf("%s %v accepted (Value=%v)", tc.name, tc.value, r.Value)
			}
		})
	}
}

func TestMoveAxisRejectsBadAxis(t *testing.T) {
	for _, axis := range []string{"99", "-1", "2", "abc", ""} {
		t.Run("axis="+axis, func(t *testing.T) {
			c := newTestClient(t)
			r, err := c.put("moveaxis", url.Values{"Axis": {axis}, "Rate": {"1.0"}}, nil)
			if err != nil {
				t.Fatalf("moveaxis Axis=%q: %v", axis, err)
			}
			if r.ErrorNumber == 0 {
				t.Errorf("moveaxis Axis=%q accepted", axis)
			}
		})
	}
}

func TestPulseGuideRejectsBadDirection(t *testing.T) {
	for _, dir := range []string{"-1", "4", "9", "abc"} {
		t.Run("direction="+dir, func(t *testing.T) {
			c := newTestClient(t)
			r, err := c.put("pulseguide", url.Values{"Direction": {dir}, "Duration": {"100"}}, nil)
			if err != nil {
				t.Fatalf("pulseguide Direction=%q: %v", dir, err)
			}
			if r.ErrorNumber == 0 {
				t.Errorf("pulseguide Direction=%q accepted", dir)
			}
		})
	}
}

func TestClientTransactionIDEcho(t *testing.T) {
	c := newTestClient(t)
	for _, cid := range []uint32{7, 123, 0, 4294967295} {
		r, err := c.get("description", url.Values{"ClientTransactionID": {strconv.FormatUint(uint64(cid), 10)}})
		if err != nil {
			t.Fatalf("description: %v", err)
		}
		if r.ClientTransactionID != cid {
			t.Errorf("ClientTransactionID echo %d -> %d", cid, r.ClientTransactionID)
		}
	}

	r1, _ := c.get("description", nil)
	r2, _ := c.get("description", nil)
	if r2.ServerTransactionID <= r1.ServerTransactionID {
		t.Error("ServerTransactionID is not monotonic")
	}
}

func TestClientIDDistinctFromTransactionID(t *testing.T) {
	c := newTestClient(t)
	r, err := c.get("description", url.Values{"ClientID": {"7"}, "ClientTransactionID": {"123"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.ClientTransactionID != 123 {
		t.Errorf("ClientTransactionID = %d, expected 123", r.ClientTransactionID)
	}
	if r.ClientTransactionID == 7 {
		t.Error("ClientTransactionID must not equal ClientID")
	}
}

func TestUnimplementedCapabilityHonest(t *testing.T) {
	c := newTestClient(t)

	r, err := c.put("synctocoordinates",
		url.Values{"RightAscension": {"1.0"}, "Declination": {"0.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("Sync declared absent but returned success")
	}

	r, err = c.put("setpark", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("SetPark declared absent but returned success")
	}
}

func TestWrongMethodAndRoute(t *testing.T) {
	c := newTestClient(t)
	// GET on a PUT-only route (park) must not execute the operation.
	resp, err := c.http.Get(c.telescopeURL("park") + "?ClientID=1&ClientTransactionID=1")
	if err != nil {
		return // network error is also acceptable (mirrors pytest raises)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		t.Errorf("GET /park returned %d (expected an error)", resp.StatusCode)
	}
}

func TestErrorThenRecover(t *testing.T) {
	c := newTestClient(t)
	r, err := c.put("synctocoordinates", url.Values{"RightAscension": {"1"}, "Declination": {"0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Fatal("sync should fail")
	}
	if mustString(t, c, "description") == "" {
		t.Error("channel not healthy after the error")
	}
}

func TestValidEdges(t *testing.T) {
	c := newTestClient(t)
	if v, err := c.getValue("trackingrates", nil); err != nil {
		t.Fatal(err)
	} else if s := v.([]any); len(s) != 3 || asInt(s[0]) != 0 || asInt(s[1]) != 1 || asInt(s[2]) != 2 {
		t.Errorf("trackingrates = %v, expected [0 1 2]", v)
	}
	if mustFloat(t, c, "guideraterightascension") <= 0 {
		t.Error("guideraterightascension not positive")
	}
	if mustFloat(t, c, "guideratedeclination") <= 0 {
		t.Error("guideratedeclination not positive")
	}
}
