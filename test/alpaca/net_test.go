//go:build integration

// NET — Network, concurrent clients and recovery.
package alpaca

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestControlRestObserveAlpaca(t *testing.T) {
	t.Skip("requires firmware /api/status (port 80); the Alpaca equivalent is devicestate")
}

func TestSlowClientDoesNotBlockService(t *testing.T) {
	c := newTestClient(t)

	host := strings.TrimPrefix(baseURL(), "http://")

	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Skipf("could not open raw socket: %v", err)
	}
	defer conn.Close()

	// send only a prefix of a request, keeping it incomplete
	fmt.Fprintf(conn, "PUT /api/v1/telescope/0/slewto")

	// the service stays healthy for a new client while the stalled one lingers
	if mustString(t, c, "description") == "" {
		t.Error("description empty")
	}
	if mustBool(t, c, "slewing") {
		t.Error("unexpected slewing")
	}
}
