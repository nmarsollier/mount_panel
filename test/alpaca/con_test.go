//go:build integration

// CON — Connection, identity and availability.
package alpaca

import (
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestIdentityVersionCapabilities(t *testing.T) {
	c := newTestClient(t)

	if got := mustInt(t, c, "interfaceversion"); got != interfaceVersion {
		t.Errorf("interfaceversion = %d, expected %d", got, interfaceVersion)
	}
	if got := mustInt(t, c, "alignmentmode"); got != alignmentMode {
		t.Errorf("alignmentmode = %d, expected %d", got, alignmentMode)
	}
	if got := mustInt(t, c, "equatorialsystem"); got != equatorialSystem {
		t.Errorf("equatorialsystem = %d, expected %d", got, equatorialSystem)
	}

	// required capabilities: true (not skipped via CanX=false)
	for _, prop := range requiredCapabilities {
		if !mustBool(t, c, prop) {
			t.Errorf("required capability %s is false", prop)
		}
	}

	// move axis requires Axis (0/1)
	for _, axis := range []int{0, 1} {
		v, err := c.boolValue("canmoveaxis", url.Values{"Axis": {strconv.Itoa(axis)}})
		if err != nil {
			t.Fatalf("canmoveaxis(Axis=%d): %v", axis, err)
		}
		if !v {
			t.Errorf("canmoveaxis(Axis=%d) false", axis)
		}
	}

	// absent optionals: not advertised as present
	if mustBool(t, c, "canslewaltaz") {
		t.Error("canslewaltaz should be false")
	}
	if mustBool(t, c, "cansync") {
		t.Error("cansync should be false")
	}
	if mustBool(t, c, "cansetpark") {
		t.Error("cansetpark should be false")
	}
	if mustBool(t, c, "cansetpierside") {
		t.Error("cansetpierside should be false")
	}

	// non-empty identity
	for _, name := range []string{"name", "description", "driverinfo", "driverversion"} {
		if mustString(t, c, name) == "" {
			t.Errorf("%s empty", name)
		}
	}

	if v, err := c.getValue("trackingrates", nil); err != nil {
		t.Fatal(err)
	} else if s := v.([]any); len(s) != 3 || asInt(s[0]) != 0 || asInt(s[1]) != 1 || asInt(s[2]) != 2 {
		t.Errorf("trackingrates = %v, expected [0 1 2]", v)
	}
}

func TestSlewSettleTimeHonest(t *testing.T) {
	c := newTestClient(t)
	r, err := c.put("slewsettletime", url.Values{"SlewSettleTime": {"5"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.ErrorNumber == 0 {
		t.Error("SlewSettleTime PUT accepted though unsupported")
	}
}

func TestConnectDisconnect(t *testing.T) {
	c := newTestClient(t)
	c.putOK(t, "connected", url.Values{"Connected": {"true"}}, nil)
	if !mustBool(t, c, "connected") {
		t.Error("connected not true after connect")
	}

	before := currentRadec(t, c)
	c.putOK(t, "connected", url.Values{"Connected": {"false"}}, nil)
	if mustBool(t, c, "connected") {
		t.Error("connected not false after disconnect")
	}

	after := currentRadec(t, c)
	if before.slewing || after.slewing {
		t.Error("connect/disconnect started motion")
	}
	if before.tracking || after.tracking {
		t.Error("connect/disconnect started tracking")
	}

	c.putOK(t, "connected", url.Values{"Connected": {"true"}}, nil)
	if !mustBool(t, c, "connected") {
		t.Error("connected not true after reconnect")
	}
}

func TestConnectIdempotent(t *testing.T) {
	c := newTestClient(t)
	for i := 0; i < 5; i++ {
		c.putOK(t, "connected", url.Values{"Connected": {"true"}}, nil)
		if !mustBool(t, c, "connected") {
			t.Fatal("connected not true")
		}
	}
	for i := 0; i < 3; i++ {
		c.putOK(t, "connected", url.Values{"Connected": {"false"}}, nil)
		if mustBool(t, c, "connected") {
			t.Fatal("connected not false")
		}
	}
	c.putOK(t, "connected", url.Values{"Connected": {"true"}}, nil)
	after := currentRadec(t, c)
	if after.slewing || after.tracking {
		t.Error("repeated connection started motion")
	}
}

func TestQueryDuringMotion(t *testing.T) {
	c := newTestClient(t)
	setupSiteTime(t, c)
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {"6.0"}, "Declination": {"-60.0"}}, nil)
	waitSlewing(t, c)

	for i := 0; i < 5; i++ {
		if !mustBool(t, c, "slewing") {
			t.Fatal("the query announced a false end of slew")
		}
		if _, err := c.floatValue("rightascension", nil); err != nil {
			t.Fatal(err)
		}
	}
	c.putOK(t, "abortslew", nil, nil)
	waitNotSlewing(t, c)
	if mustBool(t, c, "slewing") {
		t.Error("slewing still true after abortslew")
	}
}

func TestDiscoverFromNewClient(t *testing.T) {
	c := newTestClient(t)

	// baseURL() is http://host:port — strip the scheme and port for discovery.
	hostPort := strings.TrimPrefix(baseURL(), "http://")
	host := hostPort
	expectedPort := 11111
	if i := strings.LastIndex(hostPort, ":"); i >= 0 {
		host = hostPort[:i]
		if p, err := strconv.Atoi(hostPort[i+1:]); err == nil {
			expectedPort = p
		}
	}

	disc := discover(t, host)
	port, ok := disc["AlpacaPort"].(float64)
	if !ok {
		t.Fatalf("discovery without AlpacaPort: %v", disc)
	}
	if int(port) != expectedPort {
		t.Fatalf("unexpected AlpacaPort: %v", port)
	}

	mgmt, err := c.management("configureddevices")
	if err != nil {
		t.Fatal(err)
	}
	devs, ok := mgmt.Value.([]any)
	if !ok || len(devs) != 1 {
		t.Fatalf("expected 1 device, got %v", mgmt.Value)
	}
	dev, ok := devs[0].(map[string]any)
	if !ok {
		t.Fatalf("device is not an object: %v", devs[0])
	}
	if dev["DeviceType"] != "Telescope" {
		t.Errorf("DeviceType != Telescope: %v", dev["DeviceType"])
	}
	if asInt(dev["DeviceNumber"]) != deviceNumber {
		t.Errorf("DeviceNumber != %d", deviceNumber)
	}
	if disc["ServerName"] != dev["DeviceName"] {
		t.Errorf("ServerName (%v) != DeviceName (%v)", disc["ServerName"], dev["DeviceName"])
	}
}

func TestRestStatusParseable(t *testing.T) {
	t.Skip("requires firmware /api/status (port 80); the Alpaca equivalent is devicestate, covered in v4_test.go")
}
