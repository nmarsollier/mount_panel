package eafproto

import (
	"encoding/json"
	"testing"
)

// TestBuildRequestLocksTheWireShape pins the request envelope to the firmware
// contract (nrm_bft_eaf/main/proto): the exact top-level keys and which verb
// field carries the operation name.
func TestBuildRequestLocksTheWireShape(t *testing.T) {
	cases := []struct {
		name     string
		reqType  string
		target   string
		verb     string
		params   map[string]any
		wantKeys []string
		wantVerb string // "action" or "op"; "" = no verb field
	}{
		{
			name:     "capabilities",
			reqType:  "capabilities",
			wantKeys: []string{"proto", "type", "session", "boot_id"},
		},
		{
			name:     "state focuser",
			reqType:  "state",
			target:   "focuser",
			wantKeys: []string{"proto", "type", "session", "boot_id", "target"},
		},
		{
			name:     "action move",
			reqType:  "action",
			target:   "focuser",
			verb:     "move",
			params:   map[string]any{"delta_steps": int64(320)},
			wantKeys: []string{"proto", "type", "session", "boot_id", "target", "action", "params"},
			wantVerb: "action",
		},
		{
			name:     "control fan mode",
			reqType:  "control",
			target:   "fan",
			verb:     "mode",
			params:   map[string]any{"mode": "normal"},
			wantKeys: []string{"proto", "type", "session", "boot_id", "target", "op", "params"},
			wantVerb: "op",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := buildRequest(tc.reqType, tc.target, tc.verb, "g-sess", "a1b2c3d4", tc.params)

			b, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			for _, k := range tc.wantKeys {
				if _, ok := got[k]; !ok {
					t.Errorf("missing key %q in %s", k, b)
				}
			}
			if got["proto"] != float64(2) {
				t.Errorf("proto = %v, want 2", got["proto"])
			}
			if got["session"] != "g-sess" || got["boot_id"] != "a1b2c3d4" {
				t.Errorf("identity not echoed: session=%v boot_id=%v", got["session"], got["boot_id"])
			}
			if tc.wantVerb == "" {
				if _, ok := got["action"]; ok {
					t.Errorf("unexpected action field: %s", b)
				}
				if _, ok := got["op"]; ok {
					t.Errorf("unexpected op field: %s", b)
				}
			} else if got[tc.wantVerb] != tc.verb {
				t.Errorf("%s = %v, want %q", tc.wantVerb, got[tc.wantVerb], tc.verb)
			}
		})
	}
}

// TestParseErrorAcceptsBothForms covers the object form the accessory replies
// with and the string form the transport injects on a lost link.
func TestParseErrorAcceptsBothForms(t *testing.T) {
	obj := parseError(json.RawMessage(`{"code":"BUSY","message":"focuser moving","field":null}`))
	if err, ok := obj.(*Error); !ok || err.Code != "BUSY" || err.Message != "focuser moving" {
		t.Fatalf("object form: %#v", obj)
	}

	str := parseError(json.RawMessage(`"serial request timed out"`))
	if err, ok := str.(*Error); !ok || err.Message != "serial request timed out" {
		t.Fatalf("string form: %#v", str)
	}
}
