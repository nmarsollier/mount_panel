package alpaca_mount

import (
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/mount"
)

// putTracking enables (sidereal) or disables tracking.
func putTracking(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	enabled, ok := server.ParamBool(c, "Tracking")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Tracking", cid, stx)
		return
	}
	if err := mount.Get().SetTracking(enabled); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putTrackingRate sets the tracking rate (0 sidereal, 1 lunar, 2 solar).
func putTrackingRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	rate, _ := server.ParamInt(c, "TrackingRate")
	if err := mount.Get().SetTrackingRate(rate); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putUTCDate sets the mount clock from an ISO 8601 string.
func putUTCDate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	val := server.Param(c, "UTCDate")
	if val == "" {
		server.WriteError(c, server.ErrInvalidValue, "Missing UTCDate", cid, stx)
		return
	}
	if err := mount.Get().SetUTCDate(val); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSiteLatitude updates and persists the site latitude.
func putSiteLatitude(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	v, ok := server.ParamFloat(c, "SiteLatitude")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing SiteLatitude", cid, stx)
		return
	}
	if err := mount.Get().SetSiteLatitude(v); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSiteLongitude updates and persists the site longitude.
func putSiteLongitude(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	v, ok := server.ParamFloat(c, "SiteLongitude")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing SiteLongitude", cid, stx)
		return
	}
	if err := mount.Get().SetSiteLongitude(v); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSiteElevation updates and persists the site elevation.
func putSiteElevation(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	v, ok := server.ParamInt(c, "SiteElevation")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing SiteElevation", cid, stx)
		return
	}
	if err := mount.Get().SetSiteElevation(v); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putTargetRightAscension stores the target right ascension for later slew.
func putTargetRightAscension(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	v, ok := server.ParamFloat(c, "TargetRightAscension")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing TargetRightAscension", cid, stx)
		return
	}
	if err := mount.Get().SetTargetRA(v); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putTargetDeclination stores the target declination for later slew.
func putTargetDeclination(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	v, ok := server.ParamFloat(c, "TargetDeclination")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing TargetDeclination", cid, stx)
		return
	}
	if err := mount.Get().SetTargetDec(v); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSideOfPier is read-only on this mount.
func putSideOfPier(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "SideOfPier is read-only", cid, stx)
}

// putSlewSettleTime is not supported on this mount.
func putSlewSettleTime(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "SlewSettleTime is not supported", cid, stx)
}

// putDeclinationRate is managed internally.
func putDeclinationRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Not implemented", cid, stx)
}

// putRightAscensionRate is managed internally.
func putRightAscensionRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Not implemented", cid, stx)
}

// putGuideRateDeclination sets the DEC guide rate in deg/s (0-1).
func putGuideRateDeclination(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	val, present := server.ParamRaw(c, "GuideRateDeclination")
	if !present {
		server.WriteError(c, server.ErrInvalidValue, "Missing GuideRateDeclination", cid, stx)
		return
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil || math.IsNaN(rate) || math.IsInf(rate, 0) {
		server.WriteError(c, server.ErrInvalidValue, "Invalid GuideRateDeclination", cid, stx)
		return
	}
	if rate < 0 || rate > 1 {
		server.WriteError(c, server.ErrInvalidValue, "GuideRateDeclination out of range (0-1 deg/s)", cid, stx)
		return
	}
	if err := mount.Get().SetGuideRateDec(rate); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putGuideRateRightAscension sets the RA guide rate in deg/s (0-1).
func putGuideRateRightAscension(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	val, present := server.ParamRaw(c, "GuideRateRightAscension")
	if !present {
		server.WriteError(c, server.ErrInvalidValue, "Missing GuideRateRightAscension", cid, stx)
		return
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil || math.IsNaN(rate) || math.IsInf(rate, 0) {
		server.WriteError(c, server.ErrInvalidValue, "Invalid GuideRateRightAscension", cid, stx)
		return
	}
	if rate < 0 || rate > 1 {
		server.WriteError(c, server.ErrInvalidValue, "GuideRateRightAscension out of range (0-1 deg/s)", cid, stx)
		return
	}
	if err := mount.Get().SetGuideRateRA(rate); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}
