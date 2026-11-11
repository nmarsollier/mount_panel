package mount

import (
	"testing"

	"github.com/nmarsollier/mount_panel/internal/mountproto"
)

func TestPierSideStateTranslation(t *testing.T) {
	cases := map[string]int{
		"east":       0,
		"west":       1,
		"":           0, // unknown falls back to east (ASCOM pierEast)
		"unexpected": 0,
	}
	for side, want := range cases {
		got := pierSideState(&mountproto.StateResponse{PierSide: side})
		if got != want {
			t.Errorf("pierSideState(%q) = %d, want %d", side, got, want)
		}
	}
}

func TestSideOfPierGatesOnPositionUncertain(t *testing.T) {
	cases := []struct {
		name         string
		meridianFlip bool
		uncertain    bool
		pierSide     string
		want         int
	}{
		{"fork is unknown", false, false, "east", pierUnknown},
		{"uncertain position is unknown", true, true, "east", pierUnknown},
		{"east when settled", true, false, "east", 0},
		{"west when settled", true, false, "west", 1},
	}
	for _, c := range cases {
		got := sideOfPier(&mountproto.StateResponse{
			PierSide:          c.pierSide,
			PositionUncertain: c.uncertain,
		}, c.meridianFlip)
		if got != c.want {
			t.Errorf("%s: sideOfPier = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestDestinationSideFromResponseRejectsUnknown(t *testing.T) {
	cases := map[string]int{
		"east": 0,
		"west": 1,
	}
	for side, want := range cases {
		got, err := destinationSideFromResponse(side)
		if err != nil {
			t.Errorf("destinationSideFromResponse(%q) error: %v", side, err)
			continue
		}
		if got != want {
			t.Errorf("destinationSideFromResponse(%q) = %d, want %d", side, got, want)
		}
	}
	for _, side := range []string{"", "unknown", "invalid", "north"} {
		got, err := destinationSideFromResponse(side)
		if err == nil {
			t.Errorf("destinationSideFromResponse(%q) = %d, want an error", side, got)
		}
	}
}

func TestDestinationSideOfPierForkDoesNotQuery(t *testing.T) {
	// A fork (bft) reports meridian_flip=false and has no side of pier. The
	// gateway must return pierUnknown without sending the destination_side_of_pier
	// query — the connection is nil here, so any query would fail with
	// "mount not connected".
	m := &Mount{meridianFlip: false}
	got, err := m.DestinationSideOfPier(12.5, -35.0)
	if err != nil {
		t.Fatalf("fork DestinationSideOfPier returned error: %v", err)
	}
	if got != pierUnknown {
		t.Errorf("fork DestinationSideOfPier = %d, want pierUnknown(-1)", got)
	}
}
