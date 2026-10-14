package mount

import "github.com/nmarsollier/mount_panel/internal/mountproto"

// UTCDate returns the current UTC date as an ISO 8601 string, as reported by
// the mount (the device is the time source).
func (m *Mount) UTCDate() (string, error) {
	cfg, err := m.fetchConfig()
	if err != nil {
		return "", err
	}
	return cfg.Utc, nil
}

// SetUTCDate writes an explicit clock value to the mount and puts it into
// manual mode, disabling the periodic time sync until SyncTimeToHost is called.
func (m *Mount) SetUTCDate(iso string) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	if _, err := c.ConfigSet(mountproto.ConfigRequest{Utc: &iso}); err != nil {
		return err
	}
	m.mu.Lock()
	m.manualUTC = &iso
	m.mu.Unlock()
	return nil
}

// SyncTimeToHost clears any manual override and re-syncs the mount clock to
// the host, restoring the automatic time authority.
func (m *Mount) SyncTimeToHost() error {
	m.mu.Lock()
	m.manualUTC = nil
	m.mu.Unlock()
	return m.provideTime()
}
