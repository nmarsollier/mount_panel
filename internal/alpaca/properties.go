package alpaca

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

func getAlignmentMode(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, 2, cid, stx)
}

func getAltitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().Altitude()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getApertureArea(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, apertureArea, cid, stx)
}

func getApertureDiameter(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, apertureDiameter, cid, stx)
}

func getAtHome(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().AtHome()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeBool(c, v, cid, stx)
}

func getAtPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().AtPark()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeBool(c, v, cid, stx)
}

func getAzimuth(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().Azimuth()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().Declination()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getDeclinationRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, 0, cid, stx)
}

func getEquatorialSystem(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, 1, cid, stx)
}

func getFocalLength(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, focalLength, cid, stx)
}

func getGuideRateDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().GuideRateDec()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getGuideRateRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().GuideRateRA()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getIsPulseGuiding(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().IsPulseGuiding()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeBool(c, v, cid, stx)
}

func getRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().RightAscension()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getRightAscensionRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, 0, cid, stx)
}

func getSideOfPier(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().SideOfPier()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeInt(c, v, cid, stx)
}

func getSiderealTime(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().SiderealTime()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getSiteElevation(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().SiteElevation()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeInt(c, v, cid, stx)
}

func getSiteLatitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().SiteLatitude()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getSiteLongitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().SiteLongitude()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeFloat(c, v, cid, stx)
}

func getSlewing(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().Slewing()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeBool(c, v, cid, stx)
}

func getSlewSettleTime(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, 0, cid, stx)
}

func getTargetDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().TargetDec(), cid, stx)
}

func getTargetRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().TargetRA(), cid, stx)
}

func getTracking(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().Tracking()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeBool(c, v, cid, stx)
}

func getTrackingRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, mount.Get().SelectedTrackingRate(), cid, stx)
}

func getTrackingRates(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeValue(c, rawJSON([]int{0, 1, 2}), cid, stx)
}

func getUTCDate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	v, err := mount.Get().UTCDate()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeString(c, v, cid, stx)
}

func getAxisRates(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	m := mount.Get()

	// Axis is 0 (RA) or 1 (DEC); this mount uses the same range for both.
	_, _ = paramInt(c, "Axis")

	lo, err := m.MinSlewSpeed()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	hi, err := m.MaxSlewSpeed()
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeValue(c, []byte(`[{"Minimum":`+strconv.FormatFloat(lo, 'f', 1, 64)+`,"Maximum":`+strconv.FormatFloat(hi, 'f', 1, 64)+`}]`), cid, stx)
}

func getDoesRefraction(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getDestinationSideOfPier(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	ra, okRa := paramFloat(c, "RightAscension")
	dec, okDec := paramFloat(c, "Declination")
	if !okRa || !okDec {
		writeError(c, errInvalidValue, "Missing RightAscension or Declination", cid, stx)
		return
	}
	pier, err := mount.Get().DestinationSideOfPier(ra, dec)
	if err != nil {
		writeError(c, errNotConnected, err.Error(), cid, stx)
		return
	}
	writeInt(c, pier, cid, stx)
}
