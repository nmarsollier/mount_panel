//go:build integration

// PRK — Park, unpark and session close.
package alpaca

import (
	"net/url"
	"testing"
)

func TestParkFromObservation(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "tracking", url.Values{"Tracking": {"true"}}, nil)
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("atpark not true after park")
	}
	if mustBool(t, c, "tracking") {
		t.Error("tracking still active after park")
	}
}

func TestParkAgain(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("atpark not true after park")
	}
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("idempotent park failed")
	}
}

func TestParkedBlocksMotion(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "park", nil, nil)
	r, err := c.put("slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("GOTO while parked was accepted")
	}
}

func TestParkedRejectsAllCommands(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "park", nil, nil)

	rejected := []struct {
		name string
		form url.Values
	}{
		{"slewtocoordinatesasync", url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}},
		{"moveaxis", url.Values{"Axis": {"0"}, "Rate": {"1.0"}}},
		{"moveaxis", url.Values{"Axis": {"1"}, "Rate": {"1.0"}}},
		{"findhome", nil},
		{"pulseguide", url.Values{"Direction": {"0"}, "Duration": {"100"}}},
		{"tracking", url.Values{"Tracking": {"true"}}},
		{"slewtotargetasync", nil},
	}
	for _, rc := range rejected {
		r, err := c.put(rc.name, rc.form, nil)
		if err != nil {
			t.Fatalf("%s: %v", rc.name, err)
		}
		if r.ErrorNumber == 0 {
			t.Errorf("%s accepted while parked", rc.name)
		}
	}

	c.putOK(t, "unpark", nil, nil)
	if mustBool(t, c, "atpark") {
		t.Error("atpark not false after unpark")
	}
}

func TestUnpark(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "park", nil, nil)
	if !mustBool(t, c, "atpark") {
		t.Error("atpark not true")
	}
	c.putOK(t, "unpark", nil, nil)
	if mustBool(t, c, "atpark") {
		t.Error("atpark not false")
	}

	r, err := c.put("slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber != 0 {
		t.Error("after unpark the GOTO was rejected")
	}
	c.putOK(t, "abortslew", nil, nil)
}
