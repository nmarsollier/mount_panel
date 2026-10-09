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
	writeFloat(c, mount.Get().Altitude(), cid, stx)
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
	writeBool(c, mount.Get().AtHome(), cid, stx)
}

func getAtPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, mount.Get().AtPark(), cid, stx)
}

func getAzimuth(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().Azimuth(), cid, stx)
}

func getDeclination(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().Declination(), cid, stx)
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
	writeFloat(c, mount.Get().GuideRateDec(), cid, stx)
}

func getGuideRateRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().GuideRateRA(), cid, stx)
}

func getIsPulseGuiding(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, mount.Get().IsPulseGuiding(), cid, stx)
}

func getRightAscension(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().RightAscension(), cid, stx)
}

func getRightAscensionRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, 0, cid, stx)
}

func getSideOfPier(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, mount.Get().SideOfPier(), cid, stx)
}

func getSiderealTime(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().SiderealTime(), cid, stx)
}

func getSiteElevation(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeInt(c, mount.Get().SiteElevation(), cid, stx)
}

func getSiteLatitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().SiteLatitude(), cid, stx)
}

func getSiteLongitude(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeFloat(c, mount.Get().SiteLongitude(), cid, stx)
}

func getSlewing(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, mount.Get().Slewing(), cid, stx)
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
	writeBool(c, mount.Get().Tracking(), cid, stx)
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
	writeString(c, mount.Get().UTCDate(), cid, stx)
}

func getAxisRates(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	m := mount.Get()

	// Axis is 0 (RA) or 1 (DEC); this mount uses the same range for both.
	_, _ = paramInt(c, "Axis")

	lo := strconv.FormatFloat(m.MinSlewSpeed(), 'f', 1, 64)
	hi := strconv.FormatFloat(m.MaxSlewSpeed(), 'f', 1, 64)
	writeValue(c, []byte(`[{"Minimum":`+lo+`,"Maximum":`+hi+`}]`), cid, stx)
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
	writeInt(c, mount.Get().DestinationSideOfPier(ra, dec), cid, stx)
}
