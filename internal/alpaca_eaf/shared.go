package alpaca_eaf

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
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
