package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

// Stop halts everything (slew, tracking and manual motion).
func (m *Mount) Stop() error {
	return m.action(mountproto.ActionRequest{Action: "stop", Scope: "all"})
}
