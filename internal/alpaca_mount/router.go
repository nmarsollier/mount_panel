package alpaca_mount

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_eaf"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"github.com/nmarsollier/nrm_ha_go/internal/tools/env"
)

const telescope = "/api/v1/telescope/0"

// Start builds the router, registers every Alpaca endpoint, starts the UDP
// discovery listener and runs the HTTP server.
func Start() {
	engine := server.Router()

	initManagement(engine)
	initTelescope(engine)
	alpaca_eaf.Init(engine)
	startDiscovery()

	engine.Run(fmt.Sprintf(":%d", env.Get().Port))
}

// initManagement registers the management (discovery) endpoints at both the
// versioned and unversioned paths.
func initManagement(engine *gin.Engine) {
	engine.GET("/management/apiversions", getAPIVersions)
	engine.GET("/management/v1/apiversions", getAPIVersions)
	engine.GET("/management/description", getManagementDescription)
	engine.GET("/management/v1/description", getManagementDescription)
	engine.GET("/management/configureddevices", getConfiguredDevices)
	engine.GET("/management/v1/configureddevices", getConfiguredDevices)
}

// initTelescope registers the ITelescopeV4 endpoints for device 0.
func initTelescope(engine *gin.Engine) {
	g := engine.Group(telescope)

	// Common device endpoints.
	g.GET("/connected", getConnected)
	g.PUT("/connected", putConnected)
	g.GET("/connecting", getConnecting)
	g.GET("/description", getDescription)
	g.GET("/driverinfo", getDriverInfo)
	g.GET("/driverversion", getDriverVersion)
	g.GET("/interfaceversion", getInterfaceVersion)
	g.GET("/name", getName)
	g.GET("/supportedactions", getSupportedActions)
	g.GET("/devicestate", getDeviceState)

	// Capabilities.
	g.GET("/canfindhome", getCanFindHome)
	g.GET("/canpark", getCanPark)
	g.GET("/canpulseguide", getCanPulseGuide)
	g.GET("/cansetdeclinationrate", getCanSetDeclinationRate)
	g.GET("/cansetguiderates", getCanSetGuideRates)
	g.GET("/cansetpark", getCanSetPark)
	g.GET("/cansetpierside", getCanSetPierSide)
	g.GET("/cansetrightascensionrate", getCanSetRightAscensionRate)
	g.GET("/cansettracking", getCanSetTracking)
	g.GET("/canslew", getCanSlew)
	g.GET("/canslewaltaz", getCanSlewAltAz)
	g.GET("/canslewaltazasync", getCanSlewAltAzAsync)
	g.GET("/canslewasync", getCanSlewAsync)
	g.GET("/cansync", getCanSync)
	g.GET("/cansyncaltaz", getCanSyncAltAz)
	g.GET("/canunpark", getCanUnpark)
	g.GET("/canmoveaxis", getCanMoveAxis)

	// Telescope properties GET.
	g.GET("/alignmentmode", getAlignmentMode)
	g.GET("/altitude", getAltitude)
	g.GET("/aperturearea", getApertureArea)
	g.GET("/aperturediameter", getApertureDiameter)
	g.GET("/athome", getAtHome)
	g.GET("/atpark", getAtPark)
	g.GET("/azimuth", getAzimuth)
	g.GET("/declination", getDeclination)
	g.GET("/declinationrate", getDeclinationRate)
	g.GET("/equatorialsystem", getEquatorialSystem)
	g.GET("/focallength", getFocalLength)
	g.GET("/guideratedeclination", getGuideRateDeclination)
	g.GET("/guideraterightascension", getGuideRateRightAscension)
	g.GET("/ispulseguiding", getIsPulseGuiding)
	g.GET("/rightascension", getRightAscension)
	g.GET("/rightascensionrate", getRightAscensionRate)
	g.GET("/sideofpier", getSideOfPier)
	g.GET("/siderealtime", getSiderealTime)
	g.GET("/siteelevation", getSiteElevation)
	g.GET("/sitelatitude", getSiteLatitude)
	g.GET("/sitelongitude", getSiteLongitude)
	g.GET("/slewing", getSlewing)
	g.GET("/slewsettletime", getSlewSettleTime)
	g.GET("/targetdeclination", getTargetDeclination)
	g.GET("/targetrightascension", getTargetRightAscension)
	g.GET("/tracking", getTracking)
	g.GET("/trackingrate", getTrackingRate)
	g.GET("/trackingrates", getTrackingRates)
	g.GET("/utcdate", getUTCDate)

	// Optional property GETs.
	g.GET("/axisrates", getAxisRates)
	g.GET("/doesrefraction", getDoesRefraction)
	g.GET("/destinationsideofpier", getDestinationSideOfPier)

	// Telescope methods PUT.
	g.PUT("/abortslew", putAbortSlew)
	g.PUT("/findhome", putFindHome)
	g.PUT("/park", putPark)
	g.PUT("/unpark", putUnpark)
	g.PUT("/setpark", putSetPark)
	g.PUT("/slewtocoordinatesasync", putSlewToCoordinatesAsync)
	g.PUT("/slewtoaltazasync", putSlewToAltAzAsync)
	g.PUT("/synctocoordinates", putSyncToCoordinates)
	g.PUT("/synctoaltaz", putSyncToAltAz)
	g.PUT("/slewtotargetasync", putSlewToTargetAsync)
	g.PUT("/synctotarget", putSyncToTarget)
	g.PUT("/pulseguide", putPulseGuide)
	g.PUT("/moveaxis", putMoveAxis)

	// Property setters.
	g.PUT("/tracking", putTracking)
	g.PUT("/trackingrate", putTrackingRate)
	g.PUT("/utcdate", putUTCDate)
	g.PUT("/sitelatitude", putSiteLatitude)
	g.PUT("/sitelongitude", putSiteLongitude)
	g.PUT("/siteelevation", putSiteElevation)
	g.PUT("/targetrightascension", putTargetRightAscension)
	g.PUT("/targetdeclination", putTargetDeclination)
	g.PUT("/sideofpier", putSideOfPier)
	g.PUT("/slewsettletime", putSlewSettleTime)
	g.PUT("/declinationrate", putDeclinationRate)
	g.PUT("/rightascensionrate", putRightAscensionRate)
	g.PUT("/guideratedeclination", putGuideRateDeclination)
	g.PUT("/guideraterightascension", putGuideRateRightAscension)
}
