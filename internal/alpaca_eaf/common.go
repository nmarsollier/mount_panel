package alpaca_eaf

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"github.com/nmarsollier/nrm_ha_go/internal/eaf"
)

// getConnected reports whether the client is logically connected.
func getConnected(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteBool(c, eaf.Get().Connected(), cid, stx)
}

// putConnected tracks the client's logical link. Connect/disconnect never moves.
func putConnected(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()

	connected, ok := server.ParamBool(c, "Connected")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Connected", cid, stx)
		return
	}
	eaf.Get().SetConnected(connected)
	server.WriteOK(c, cid, stx)
}

// putConnect establishes the logical client link.
func putConnect(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	eaf.Get().SetConnected(true)
	server.WriteOK(c, cid, stx)
}

// putDisconnect drops the logical client link.
func putDisconnect(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	eaf.Get().SetConnected(false)
	server.WriteOK(c, cid, stx)
}

// getConnecting reports whether a connection is being established (always false).
func getConnecting(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

// getDescription returns the device description string.
func getDescription(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteString(c, server.ServerDescription, cid, stx)
}

// getDriverInfo returns the driver identification string.
func getDriverInfo(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteString(c, server.DriverInfo, cid, stx)
}

// getDriverVersion returns the driver version.
func getDriverVersion(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteString(c, server.DriverVersion, cid, stx)
}

// getInterfaceVersion returns the ASCOM interface version (4).
func getInterfaceVersion(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteInt(c, server.InterfaceVersion, cid, stx)
}

// getName returns the focuser device name.
func getName(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteString(c, deviceName, cid, stx)
}

// getSupportedActions returns the list of custom action names (empty).
func getSupportedActions(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteValue(c, server.RawJSON([]string{}), cid, stx)
}

// getDeviceState returns the Focuser operational properties in one call.
func getDeviceState(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	entries := []struct {
		Name  string          `json:"Name"`
		Value json.RawMessage `json:"Value"`
	}{}

	if v, err := eaf.Get().IsMoving(); err == nil {
		entries = append(entries, struct {
			Name  string          `json:"Name"`
			Value json.RawMessage `json:"Value"`
		}{"IsMoving", server.RawBool(v)})
	}
	if v, err := eaf.Get().Position(); err == nil {
		entries = append(entries, struct {
			Name  string          `json:"Name"`
			Value json.RawMessage `json:"Value"`
		}{"Position", server.RawInt(int(v))})
	}
	if v, err := eaf.Get().Temperature(); err == nil {
		entries = append(entries, struct {
			Name  string          `json:"Name"`
			Value json.RawMessage `json:"Value"`
		}{"Temperature", server.RawFloat(v)})
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}
