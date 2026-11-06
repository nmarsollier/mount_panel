package alpaca_amb

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/amb"
)

// getHumidity returns the ambient relative humidity in percent.
func getHumidity(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := amb.Get().Humidity()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getDewPoint returns the dew point in Celsius.
func getDewPoint(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := amb.Get().DewPoint()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getTemperature returns the ambient temperature in Celsius.
func getTemperature(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := amb.Get().AmbientTemperature()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getAveragePeriod returns the averaging window in hours.
func getAveragePeriod(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, amb.Get().AveragePeriod(), cid, stx)
}

// putAveragePeriod sets the averaging window in hours (0..1).
func putAveragePeriod(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, ok := server.ParamFloat(c, "AveragePeriod")
	if !ok || v < 0 || v > 1 {
		server.WriteError(c, server.ErrInvalidValue, "AveragePeriod out of range (0-1 hours)", cid, stx)
		return
	}
	amb.Get().SetAveragePeriod(v)
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
