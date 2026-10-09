package mount

import "github.com/nmarsollier/nrm_ha_go/internal/proto"

// Reset reboots the ESP. It acknowledges, then restarts ~200 ms later; the
// next state poll detects the boot_id change and re-reads identity/config.
func (m *Mount) Reset() error {
	return m.action(proto.ActionRequest{Action: "reset"})
}
