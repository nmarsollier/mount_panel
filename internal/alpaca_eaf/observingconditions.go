package alpaca_eaf

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/eaf"
)

const observingConditions = "/api/v1/observingconditions/0"

// initObservingConditions registers the ObservingConditions V2 endpoints.
func initObservingConditions(engine *gin.Engine) {
	g := engine.Group(observingConditions)

	// Common device endpoints.
	g.GET("/connected", getConnected)
	g.PUT("/connected", putConnected)
	g.PUT("/connect", putConnect)
	g.PUT("/disconnect", putDisconnect)
	g.GET("/connecting", getConnecting)
	g.GET("/description", getDescription)
	g.GET("/driverinfo", getDriverInfo)
	g.GET("/driverversion", getDriverVersion)
	g.GET("/interfaceversion", getInterfaceVersionOC)
	g.GET("/name", getName)
	g.GET("/supportedactions", getSupportedActions)
	g.GET("/devicestate", getDeviceStateOC)
	g.PUT("/commandblind", putCommandBlind)
	g.PUT("/commandbool", putCommandBool)
	g.PUT("/commandstring", putCommandString)
	g.PUT("/action", putAction)

	// ObservingConditions properties.
	g.GET("/averageperiod", getAveragePeriod)
	g.PUT("/averageperiod", putAveragePeriod)
	g.GET("/cloudcover", notImplemented)
	g.GET("/dewpoint", getDewPoint)
	g.GET("/humidity", getHumidity)
	g.GET("/pressure", notImplemented)
	g.GET("/rainrate", notImplemented)
	g.GET("/skybrightness", notImplemented)
	g.GET("/skyquality", notImplemented)
	g.GET("/skytemperature", notImplemented)
	g.GET("/starfwhm", notImplemented)
	g.GET("/temperature", getTemperature)
	g.GET("/winddirection", notImplemented)
	g.GET("/windgust", notImplemented)
	g.GET("/windspeed", notImplemented)
	g.GET("/sensordescription", getSensorDescription)
	g.GET("/timesincelastupdate", getTimeSinceLastUpdate)
	g.PUT("/refresh", putRefresh)
}

// getInterfaceVersionOC returns the ASCOM interface version (2).
func getInterfaceVersionOC(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 2, cid, stx)
}

// getHumidity returns the ambient relative humidity in percent.
func getHumidity(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := eaf.Get().Humidity()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getDewPoint returns the dew point in Celsius.
func getDewPoint(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := eaf.Get().DewPoint()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getAveragePeriod returns the averaging window in hours.
func getAveragePeriod(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, eaf.Get().AveragePeriod(), cid, stx)
}

// putAveragePeriod sets the averaging window in hours (0..1).
func putAveragePeriod(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, ok := server.ParamFloat(c, "AveragePeriod")
	if !ok || v < 0 || v > 1 {
		server.WriteError(c, server.ErrInvalidValue, "AveragePeriod out of range (0-1 hours)", cid, stx)
		return
	}
	eaf.Get().SetAveragePeriod(v)
	server.WriteOK(c, cid, stx)
}

// putRefresh triggers a sensor read. The firmware does not implement the
// refresh action yet; the next state read returns a fresh sample anyway.
func putRefresh(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteOK(c, cid, stx)
}

// getSensorDescription returns a description per implemented property.
func getSensorDescription(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	switch strings.ToLower(server.Param(c, "SensorName")) {
	case "temperature":
		server.WriteString(c, "Ambient temperature (SHT31)", cid, stx)
	case "humidity":
		server.WriteString(c, "Relative humidity (SHT31)", cid, stx)
	case "dewpoint":
		server.WriteString(c, "Dew point, calculated by firmware", cid, stx)
	case "averageperiod":
		server.WriteString(c, "Averaging window, hours", cid, stx)
	case "cloudcover", "pressure", "rainrate", "skybrightness", "skyquality",
		"skytemperature", "starfwhm", "winddirection", "windgust", "windspeed":
		server.WriteError(c, server.ErrNotImplemented, "Not implemented", cid, stx)
	default:
		server.WriteError(c, server.ErrInvalidValue, "Unknown SensorName", cid, stx)
	}
}

// getTimeSinceLastUpdate returns seconds since the last valid sample.
func getTimeSinceLastUpdate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	switch strings.ToLower(server.Param(c, "SensorName")) {
	case "temperature", "humidity", "dewpoint", "averageperiod":
		server.WriteFloat(c, 0, cid, stx) // the mock always reports a fresh value
	default:
		server.WriteError(c, server.ErrInvalidValue, "Unknown SensorName", cid, stx)
	}
}

// getDeviceStateOC returns the known ObservingConditions properties.
func getDeviceStateOC(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	entries := []stateEntry{}

	if v, err := eaf.Get().Temperature(); err == nil {
		entries = append(entries, stateEntry{"Temperature", server.RawFloat(v)})
	}
	if v, err := eaf.Get().Humidity(); err == nil {
		entries = append(entries, stateEntry{"Humidity", server.RawFloat(v)})
	}
	if v, err := eaf.Get().DewPoint(); err == nil {
		entries = append(entries, stateEntry{"DewPoint", server.RawFloat(v)})
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}
