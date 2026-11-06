package alpaca_fan

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/fan"
)

// getConnected reports whether the client is logically connected.
func getConnected(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, fan.Get().Connected(), cid, stx)
}

// putConnected tracks the client's logical link.
func putConnected(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	connected, ok := server.ParamBool(c, "Connected")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Connected", cid, stx)
		return
	}
	fan.Get().SetConnected(connected)
	server.WriteOK(c, cid, stx)
}

// putConnect establishes the logical client link.
func putConnect(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	fan.Get().SetConnected(true)
	server.WriteOK(c, cid, stx)
}

// putDisconnect drops the logical client link.
func putDisconnect(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	fan.Get().SetConnected(false)
	server.WriteOK(c, cid, stx)
}

// getConnecting reports whether a connection is being established (always false).
func getConnecting(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

// getDescription returns the device description string.
func getDescription(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteString(c, deviceDescription, cid, stx)
}

// getDriverInfo returns the driver identification string.
func getDriverInfo(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteString(c, server.DriverInfo, cid, stx)
}

// getDriverVersion returns the driver version.
func getDriverVersion(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteString(c, server.DriverVersion, cid, stx)
}

// getInterfaceVersion returns the ASCOM interface version (3).
func getInterfaceVersion(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 3, cid, stx)
}

// getName returns the device name.
func getName(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteString(c, deviceName, cid, stx)
}

// getSupportedActions returns the list of custom action names (empty).
func getSupportedActions(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteValue(c, server.RawJSON([]string{}), cid, stx)
}

// getDeviceState returns the Switch operational properties.
func getDeviceState(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	entries := []stateEntry{}

	for id := 0; id < maxSwitch; id++ {
		v, err := switchValue(id)
		if err != nil {
			continue
		}
		entries = append(entries,
			stateEntry{switchStateName(id), server.RawBool(v != 0)},
			stateEntry{switchValueName(id), server.RawFloat(v)},
		)
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}
