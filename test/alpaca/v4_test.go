//go:build integration

// V4 — DeviceState ("read all") and Connecting (Telescope Interface V4).
package alpaca

import (
	"testing"
)

var operationalNames = map[string]bool{
	"Altitude": true, "AtHome": true, "AtPark": true, "Azimuth": true,
	"Declination": true, "IsPulseGuiding": true, "RightAscension": true,
	"SideOfPier": true, "Slewing": true, "Tracking": true,
}

func deviceStateByName(t *testing.T, c *Client) map[string]any {
	t.Helper()
	r, err := c.get("devicestate", nil)
	if err != nil {
		t.Fatalf("devicestate: %v", err)
	}
	if r.ErrorNumber != 0 {
		t.Fatalf("devicestate err=%d: %s", r.ErrorNumber, r.ErrorMessage)
	}
	values, ok := r.Value.([]any)
	if !ok || len(values) == 0 {
		t.Fatalf("devicestate Value empty: %v", r.Value)
	}
	out := map[string]any{}
	for _, item := range values {
		m, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("devicestate entry is not an object: %v", item)
		}
		name, _ := m["Name"].(string)
		out[name] = m["Value"]
	}
	return out
}

func TestDeviceStateOperationalNames(t *testing.T) {
	c := newTestClient(t)
	byName := deviceStateByName(t, c)
	for name := range operationalNames {
		if _, ok := byName[name]; !ok {
			t.Errorf("missing operational property: %s", name)
		}
	}
}

func TestDeviceStateValueTypes(t *testing.T) {
	c := newTestClient(t)
	byName := deviceStateByName(t, c)
	for _, name := range []string{"AtHome", "AtPark", "IsPulseGuiding", "Slewing", "Tracking"} {
		if _, ok := byName[name].(bool); !ok {
			t.Errorf("%s is not bool: %v", name, byName[name])
		}
	}
	for _, name := range []string{"Altitude", "Azimuth", "Declination", "RightAscension"} {
		switch byName[name].(type) {
		case float64, int:
		default:
			t.Errorf("%s is not numeric: %v", name, byName[name])
		}
	}
	if _, ok := byName["SideOfPier"].(float64); !ok {
		t.Errorf("SideOfPier is not numeric: %v", byName["SideOfPier"])
	}
}

func TestDeviceStateMatchesIndividual(t *testing.T) {
	c := newTestClient(t)
	byName := deviceStateByName(t, c)
	if byName["Slewing"] != mustBool(t, c, "slewing") {
		t.Error("devicestate Slewing does not match the individual endpoint")
	}
	if byName["Tracking"] != mustBool(t, c, "tracking") {
		t.Error("devicestate Tracking does not match the individual endpoint")
	}
	if byName["AtPark"] != mustBool(t, c, "atpark") {
		t.Error("devicestate AtPark does not match the individual endpoint")
	}
	if byName["AtHome"] != mustBool(t, c, "athome") {
		t.Error("devicestate AtHome does not match the individual endpoint")
	}
	if asInt(byName["SideOfPier"]) != mustInt(t, c, "sideofpier") {
		t.Error("devicestate SideOfPier does not match the individual endpoint")
	}
}

func TestConnectingIsFalse(t *testing.T) {
	c := newTestClient(t)
	if mustBool(t, c, "connecting") {
		t.Error("connecting should be false")
	}
}
