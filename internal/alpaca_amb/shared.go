package alpaca_amb

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/amb"
)

// stateEntry is a single {Name, Value} pair of a DeviceState response.
type stateEntry struct {
	Name  string          `json:"Name"`
	Value json.RawMessage `json:"Value"`
}

// notImplemented answers NotImplemented for a standard property this driver
// does not expose.
func notImplemented(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Not implemented", cid, stx)
}

// writeEnvironmentError maps an environment read failure to Alpaca: an absent
// sensor is 1024 (not implemented), a broken serial link is 1027 (not
// connected).
func writeEnvironmentError(c *gin.Context, err error, cid, stx uint32) {
	if errors.Is(err, amb.ErrUnavailable) {
		server.WriteError(c, server.ErrNotImplemented, err.Error(), cid, stx)
		return
	}
	server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
}
