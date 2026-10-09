package alpaca

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// putTracking enables (sidereal) or disables tracking.
func putTracking(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	enabled, ok := paramBool(c, "Tracking")
	if !ok {
		writeError(c, errInvalidValue, "Invalid Tracking", cid, stx)
		return
	}
	if err := mount.Get().SetTracking(enabled); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putTrackingRate sets the tracking rate (0 sidereal, 1 lunar, 2 solar).
func putTrackingRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	rate, _ := paramInt(c, "TrackingRate")
	if err := mount.Get().SetTrackingRate(rate); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putUTCDate sets the mount clock from an ISO 8601 string.
func putUTCDate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	val := param(c, "UTCDate")
	if val == "" {
		writeError(c, errInvalidValue, "Missing UTCDate", cid, stx)
		return
	}
	if err := mount.Get().SetUTCDate(val); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSiteLatitude updates and persists the site latitude.
func putSiteLatitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	v, ok := paramFloat(c, "SiteLatitude")
	if !ok {
		writeError(c, errInvalidValue, "Missing SiteLatitude", cid, stx)
		return
	}
	if err := mount.Get().SetSiteLatitude(v); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSiteLongitude updates and persists the site longitude.
func putSiteLongitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	v, ok := paramFloat(c, "SiteLongitude")
	if !ok {
		writeError(c, errInvalidValue, "Missing SiteLongitude", cid, stx)
		return
	}
	if err := mount.Get().SetSiteLongitude(v); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSiteElevation updates and persists the site elevation.
func putSiteElevation(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	v, ok := paramInt(c, "SiteElevation")
	if !ok {
		writeError(c, errInvalidValue, "Missing SiteElevation", cid, stx)
		return
	}
	if err := mount.Get().SetSiteElevation(v); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putTargetRightAscension stores the target right ascension for later slew.
func putTargetRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	v, ok := paramFloat(c, "TargetRightAscension")
	if !ok {
		writeError(c, errInvalidValue, "Missing TargetRightAscension", cid, stx)
		return
	}
	if err := mount.Get().SetTargetRA(v); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putTargetDeclination stores the target declination for later slew.
func putTargetDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	v, ok := paramFloat(c, "TargetDeclination")
	if !ok {
		writeError(c, errInvalidValue, "Missing TargetDeclination", cid, stx)
		return
	}
	if err := mount.Get().SetTargetDec(v); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSideOfPier is read-only on this mount.
func putSideOfPier(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "SideOfPier is read-only", cid, stx)
}

// putSlewSettleTime is not supported on this mount.
func putSlewSettleTime(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "SlewSettleTime is not supported", cid, stx)
}

// putDeclinationRate is managed internally.
func putDeclinationRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Not implemented", cid, stx)
}

// putRightAscensionRate is managed internally.
func putRightAscensionRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Not implemented", cid, stx)
}

// putGuideRateDeclination sets the DEC guide rate in deg/s (0-1).
func putGuideRateDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	val, present := paramRaw(c, "GuideRateDeclination")
	if !present {
		writeError(c, errInvalidValue, "Missing GuideRateDeclination", cid, stx)
		return
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil || math.IsNaN(rate) || math.IsInf(rate, 0) {
		writeError(c, errInvalidValue, "Invalid GuideRateDeclination", cid, stx)
		return
	}
	if rate < 0 || rate > 1 {
		writeError(c, errInvalidValue, "GuideRateDeclination out of range (0-1 deg/s)", cid, stx)
		return
	}
	if err := mount.Get().SetGuideRateDec(rate); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putGuideRateRightAscension sets the RA guide rate in deg/s (0-1).
func putGuideRateRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	val, present := paramRaw(c, "GuideRateRightAscension")
	if !present {
		writeError(c, errInvalidValue, "Missing GuideRateRightAscension", cid, stx)
		return
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil || math.IsNaN(rate) || math.IsInf(rate, 0) {
		writeError(c, errInvalidValue, "Invalid GuideRateRightAscension", cid, stx)
		return
	}
	if rate < 0 || rate > 1 {
		writeError(c, errInvalidValue, "GuideRateRightAscension out of range (0-1 deg/s)", cid, stx)
		return
	}
	if err := mount.Get().SetGuideRateRA(rate); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}
