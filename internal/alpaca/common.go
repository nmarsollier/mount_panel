package alpaca

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// getConnected reports whether the client is logically connected.
func getConnected(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeBool(c, mount.Get().Connected(), cid, stx)
}

// putConnected tracks the client's logical link. Connect/disconnect never moves.
func putConnected(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()

	connected, ok := paramBool(c, "Connected")
	if !ok {
		writeError(c, errInvalidValue, "Invalid Connected", cid, stx)
		return
	}
	mount.Get().SetConnected(connected)
	writeOK(c, cid, stx)
}

// getConnecting reports whether a connection is being established (always false).
func getConnecting(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeBool(c, false, cid, stx)
}

// getDescription returns the device description string.
func getDescription(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeString(c, serverDescription, cid, stx)
}

// getDriverInfo returns the driver identification string.
func getDriverInfo(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeString(c, driverInfo, cid, stx)
}

// getDriverVersion returns the driver version.
func getDriverVersion(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeString(c, driverVersion, cid, stx)
}

// getInterfaceVersion returns the ASCOM interface version (4).
func getInterfaceVersion(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeInt(c, interfaceVersion, cid, stx)
}

// getName returns the telescope device name.
func getName(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeString(c, deviceName, cid, stx)
}

// getSupportedActions returns the list of custom action names (empty).
func getSupportedActions(c *gin.Context) {
	cid := clientTransactionID(c)
	stx := nextServerTx()
	writeValue(c, rawJSON([]string{}), cid, stx)
}

// getDeviceState returns every Telescope operational property in one call, as
// a JSON array of {Name, Value} objects.
func getDeviceState(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	ds, err := mount.Get().DeviceState()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}

	entries := []struct {
		Name  string          `json:"Name"`
		Value json.RawMessage `json:"Value"`
	}{
		{"Altitude", rawFloat(ds.Altitude)},
		{"AtHome", rawBool(ds.AtHome)},
		{"AtPark", rawBool(ds.AtPark)},
		{"Azimuth", rawFloat(ds.Azimuth)},
		{"Declination", rawFloat(ds.Declination)},
		{"IsPulseGuiding", rawBool(ds.IsPulseGuiding)},
		{"RightAscension", rawFloat(ds.RightAscension)},
		{"SideOfPier", rawInt(ds.SideOfPier)},
		{"Slewing", rawBool(ds.Slewing)},
		{"Tracking", rawBool(ds.Tracking)},
	}

	value, _ := json.Marshal(entries)
	writeValue(c, value, cid, stx)
}
