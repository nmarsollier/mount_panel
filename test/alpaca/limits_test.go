//go:build integration

// LIMITS — runtime-changed limits are still enforced and reset cleanly.
//
// Skipped: setting a limit at runtime goes through the firmware's /api/limits
// debug endpoint (mount_limits_set), which has no serial-protocol equivalent
// and therefore no Alpaca-bridge endpoint. The static-limit behaviour is
// covered indirectly by got_limits_test.go (goto rejection beyond the axis
// limits).
package alpaca

import "testing"

func TestChangedDecLimitRespected(t *testing.T) {
	t.Skip("requires changing limits at runtime (firmware /api/limits); not exposed by the bridge")
}
