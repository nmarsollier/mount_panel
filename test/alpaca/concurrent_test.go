//go:build integration

// CONCURRENT — multiple real clients against the mount at once.
package alpaca

import (
	"net/url"
	"testing"
	"time"
)

func TestPollingDoesNotDisturbMotion(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	mover := newClient()
	canceller := newClient()

	mover.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, mover)

	sawSlewing := false
	for i := 0; i < 6; i++ {
		if v, err := c.boolValue("slewing", nil); err == nil && v {
			sawSlewing = true
		}
		time.Sleep(150 * time.Millisecond)
	}
	if !sawSlewing {
		t.Error("the polling client never observed the slew")
	}
	if v, err := mover.boolValue("slewing", nil); err != nil || !v {
		t.Error("the polling cancelled the motion")
	}

	canceller.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, mover)
	if v, _ := mover.boolValue("slewing", nil); v {
		t.Error("mover still slewing")
	}
}

func TestTransactionIDsCorrelate(t *testing.T) {
	c := newTestClient(t)
	a := c.nextTxn()
	b := c.nextTxn()

	ra, err := c.get("interfaceversion", url.Values{"ClientID": {"1"}, "ClientTransactionID": {uitoa(a)}})
	if err != nil {
		t.Fatal(err)
	}
	rb, err := c.get("interfaceversion", url.Values{"ClientID": {"1"}, "ClientTransactionID": {uitoa(b)}})
	if err != nil {
		t.Fatal(err)
	}
	if ra.ClientTransactionID != a {
		t.Errorf("ClientTransactionID mismatch: enviado %d, recibido %d", a, ra.ClientTransactionID)
	}
	if rb.ClientTransactionID != b {
		t.Errorf("ClientTransactionID mismatch: enviado %d, recibido %d", b, rb.ClientTransactionID)
	}
	if a == b {
		t.Error("transaction ids should be distinct")
	}
}
