//go:build integration

// MAN — Manual centering and movement (MoveAxis).
package alpaca

import (
	"net/url"
	"testing"
	"time"
)

func move(t *testing.T, c *Client, axis int, rate float64) {
	t.Helper()
	c.putOK(t, "moveaxis", url.Values{"Axis": {uitoa(uint32(axis))}, "Rate": {fs(rate)}}, nil)
}

func TestMoveAxisConstantSpeed(t *testing.T) {
	c := newTestClient(t)
	move(t, c, 0, 1.0)
	waitSlewing(t, c)
	move(t, c, 0, 0.0)
	waitNotSlewing(t, c)
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false after stopping")
	}
}

func TestMoveAxisReverse(t *testing.T) {
	t.Skip("requires reading the physical axis position (ra_deg); not observable via Alpaca")
}

func TestMoveTwoAxesStopOne(t *testing.T) {
	t.Skip("requires reading the physical axis position (dec_deg); not observable via Alpaca")
}

func TestReleaseAxisRepeatedly(t *testing.T) {
	c := newTestClient(t)
	move(t, c, 0, 1.0)
	waitSlewing(t, c)
	for i := 0; i < 3; i++ {
		move(t, c, 0, 0.0)
	}
	if mustBool(t, c, "slewing") {
		t.Error("slewing not false after releasing repeatedly")
	}
}

func TestChangeRateWithoutRelease(t *testing.T) {
	c := newTestClient(t)
	move(t, c, 0, 1.0)
	waitSlewing(t, c)
	move(t, c, 0, 3.0)
	time.Sleep(300 * time.Millisecond)
	if !mustBool(t, c, "slewing") {
		t.Error("slewing not true after raising the speed")
	}
	move(t, c, 0, 0.0)
	waitNotSlewing(t, c)
}

func TestMoveAxisLimit(t *testing.T) {
	t.Skip("requires reading the physical axis position and limits; not observable via Alpaca")
}
