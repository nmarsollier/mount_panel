package mount

import "github.com/nmarsollier/nrm_ha_go/internal/mountproto"

// Park moves the mount to the park position and marks it parked.
func (m *Mount) Park() error {
	return m.action(mountproto.ActionRequest{Action: "park"})
}

// Unpark takes the mount out of the parked state.
func (m *Mount) Unpark() error {
	return m.action(mountproto.ActionRequest{Action: "unpark"})
}

// FindHome starts the home-finding sequence (asynchronous).
func (m *Mount) FindHome() error {
	return m.action(mountproto.ActionRequest{Action: "home"})
}
