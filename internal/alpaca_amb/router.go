package alpaca_amb

import "github.com/gin-gonic/gin"

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
	g.GET("/interfaceversion", getInterfaceVersion)
	g.GET("/name", getName)
	g.GET("/supportedactions", getSupportedActions)
	g.GET("/devicestate", getDeviceState)
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
