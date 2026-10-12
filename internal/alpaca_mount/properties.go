package alpaca_mount

import (
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

func getAlignmentMode(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 2, cid, stx)
}

func getAltitude(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().Altitude()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getApertureArea(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, apertureArea, cid, stx)
}

func getApertureDiameter(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, apertureDiameter, cid, stx)
}

func getAtHome(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().AtHome()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

func getAtPark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().AtPark()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

func getAzimuth(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().Azimuth()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getDeclination(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().Declination()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getDeclinationRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, 0, cid, stx)
}

func getEquatorialSystem(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 1, cid, stx)
}

func getFocalLength(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, focalLength, cid, stx)
}

func getGuideRateDeclination(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().GuideRateDec()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getGuideRateRightAscension(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().GuideRateRA()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getIsPulseGuiding(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().IsPulseGuiding()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

func getRightAscension(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().RightAscension()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getRightAscensionRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, 0, cid, stx)
}

func getSideOfPier(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().SideOfPier()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteInt(c, v, cid, stx)
}

func getSiderealTime(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().SiderealTime()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getSiteElevation(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().SiteElevation()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteInt(c, v, cid, stx)
}

func getSiteLatitude(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().SiteLatitude()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getSiteLongitude(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().SiteLongitude()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

func getSlewing(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().Slewing()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

func getSlewSettleTime(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 0, cid, stx)
}

func getTargetDeclination(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, mount.Get().TargetDec(), cid, stx)
}

func getTargetRightAscension(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteFloat(c, mount.Get().TargetRA(), cid, stx)
}

func getTracking(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().Tracking()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

func getTrackingRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, mount.Get().SelectedTrackingRate(), cid, stx)
}

func getTrackingRates(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteValue(c, server.RawJSON([]int{0, 1, 2}), cid, stx)
}

func getUTCDate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := mount.Get().UTCDate()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteString(c, v, cid, stx)
}

func getAxisRates(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	m := mount.Get()

	// Axis is 0 (RA) or 1 (DEC); this mount uses the same range for both.
	_, _ = server.ParamInt(c, "Axis")

	lo, err := m.MinSlewSpeed()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	hi, err := m.MaxSlewSpeed()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteValue(c, []byte(`[{"Minimum":`+strconv.FormatFloat(lo, 'f', 1, 64)+`,"Maximum":`+strconv.FormatFloat(hi, 'f', 1, 64)+`}]`), cid, stx)
}

func getDoesRefraction(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getDestinationSideOfPier(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	ra, okRa := server.ParamFloat(c, "RightAscension")
	dec, okDec := server.ParamFloat(c, "Declination")
	if !okRa || !okDec {
		server.WriteError(c, server.ErrInvalidValue, "Missing RightAscension or Declination", cid, stx)
		return
	}
	pier, err := mount.Get().DestinationSideOfPier(ra, dec)
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteInt(c, pier, cid, stx)
}
