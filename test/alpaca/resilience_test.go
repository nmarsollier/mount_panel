//go:build integration

// RESILIENCE — the mount survives a lost reply without duplicating work.
package alpaca

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"
)

func rawPut(path string, form url.Values) error {
	host := strings.TrimPrefix(baseURL(), "http://")
	body := form.Encode()
	req := fmt.Sprintf("PUT %s?ClientID=1&ClientTransactionID=1 HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"Content-Type: application/x-www-form-urlencoded\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n\r\n%s", path, host, len(body), body)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(req))
	return err // never read the response — the client "lost" it
}

func TestLostResponseGotoStillExecutes(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)

	if err := rawPut("/api/v1/telescope/0/slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}); err != nil {
		t.Skipf("raw socket unavailable: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	if !mustBool(t, c, "slewing") {
		t.Error("goto did not start despite the lost reply")
	}
	waitNotSlewing(t, c)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false at the end")
	}
}

func TestMoveAxisSurvivesClientDisconnect(t *testing.T) {
	mover := newClient()
	mover.putOK(t, "moveaxis", url.Values{"Axis": {"0"}, "Rate": {"1.0"}}, nil)
	waitSlewing(t, mover)

	fresh := newClient()
	if !mustBool(t, fresh, "slewing") {
		t.Error("the motion stopped when the client left")
	}

	fresh.putOK(t, "moveaxis", url.Values{"Axis": {"0"}, "Rate": {"0.0"}}, nil)
	waitNotSlewing(t, fresh)
	if mustBool(t, fresh, "slewing") {
		t.Error("slewing not false after stopping")
	}
}
