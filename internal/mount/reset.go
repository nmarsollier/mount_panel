package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

// Reset reboots the ESP. It acknowledges, then restarts ~200 ms later; the
// next state poll detects the boot_id change and re-reads identity/config.
func (m *Mount) Reset() error {
	return m.action(mountproto.ActionRequest{Action: "reset"})
}
