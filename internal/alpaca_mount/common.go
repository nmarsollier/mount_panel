package alpaca_mount

import (
	"encoding/json"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/mount"
)

// getConnected reports whether the client is logically connected.
func getConnected(c *gin.Context) {
	cid := server.ClientTransactionID(c)
	stx := server.NextServerTx()
	server.WriteBool(c, mount.Get().Connected(), cid, stx)
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
	mount.Get().SetConnected(connected)
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

// getName returns the telescope device name.
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

// getDeviceState returns every Telescope operational property in one call, as
// a JSON array of {Name, Value} objects.
func getDeviceState(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	ds, err := mount.Get().DeviceState()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}

	entries := []struct {
		Name  string          `json:"Name"`
		Value json.RawMessage `json:"Value"`
	}{
		{"Altitude", server.RawFloat(ds.Altitude)},
		{"AtHome", server.RawBool(ds.AtHome)},
		{"AtPark", server.RawBool(ds.AtPark)},
		{"Azimuth", server.RawFloat(ds.Azimuth)},
		{"Declination", server.RawFloat(ds.Declination)},
		{"IsPulseGuiding", server.RawBool(ds.IsPulseGuiding)},
		{"RightAscension", server.RawFloat(ds.RightAscension)},
		{"SideOfPier", server.RawInt(ds.SideOfPier)},
		{"Slewing", server.RawBool(ds.Slewing)},
		{"Tracking", server.RawBool(ds.Tracking)},
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}
