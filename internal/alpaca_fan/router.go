package alpaca_fan

import "github.com/gin-gonic/gin"

const switchDevice = "/api/v1/switch/0"

// initSwitch registers the Switch V3 endpoints for device 0.
func initSwitch(engine *gin.Engine) {
	g := engine.Group(switchDevice)

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

	// Switch endpoints.
	g.GET("/maxswitch", getMaxSwitch)
	g.GET("/canwrite", getCanWrite)
	g.GET("/getswitch", getSwitch)
	g.GET("/getswitchdescription", getSwitchDescription)
	g.GET("/getswitchname", getSwitchName)
	g.PUT("/setswitchname", setSwitchName)
	g.GET("/getswitchvalue", getSwitchValue)
	g.GET("/minswitchvalue", getMinSwitchValue)
	g.GET("/maxswitchvalue", getMaxSwitchValue)
	g.GET("/switchstep", getSwitchStep)
	g.PUT("/setswitch", setSwitch)
	g.PUT("/setswitchvalue", setSwitchValue)
	g.GET("/canasync", getCanAsync)
	g.PUT("/setasync", setAsync)
	g.PUT("/setasyncvalue", setAsyncValue)
	g.GET("/statechangecomplete", getStateChangeComplete)
	g.PUT("/cancelasync", cancelAsync)
}
