package alpaca_eaf

import "github.com/gin-gonic/gin"

const focuser = "/api/v1/focuser/0"

// initFocuser registers the Focuser V4 endpoints for device 0.
func initFocuser(engine *gin.Engine) {
	g := engine.Group(focuser)

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

	// Focuser properties.
	g.GET("/absolute", getAbsolute)
	g.GET("/ismoving", getIsMoving)
	g.GET("/maxincrement", getMaxIncrement)
	g.GET("/maxstep", getMaxStep)
	g.GET("/position", getPosition)
	g.GET("/stepsize", getStepSize)
	g.GET("/tempcomp", getTempComp)
	g.GET("/tempcompavailable", getTempCompAvailable)
	g.GET("/temperature", getTemperature)

	// Focuser methods.
	g.PUT("/move", putMove)
	g.PUT("/halt", putHalt)
	g.PUT("/tempcomp", putTempComp)

	// Legacy command methods (not supported).
	g.PUT("/commandblind", putCommandBlind)
	g.PUT("/commandbool", putCommandBool)
	g.PUT("/commandstring", putCommandString)
	g.PUT("/action", putAction)
}
