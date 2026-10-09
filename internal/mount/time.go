package mount

import (
	"time"

	"github.com/nmarsollier/nrm_ha_go/internal/proto"
)

// UTCDate returns the current UTC date as an ISO 8601 string.
func (m *Mount) UTCDate() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05")
}

// SetUTCDate forwards a clock update to the mount (validation happens there).
func (m *Mount) SetUTCDate(iso string) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	_, err = c.ConfigSet(proto.ConfigRequest{Utc: &iso})
	return err
}
