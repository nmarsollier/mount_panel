package alpaca

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// notError reports whether the mount is operational (not in the error state).
func notError() bool {
	return !mount.Get().IsError()
}

func getCanFindHome(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, notError(), cid, stx)
}

func getCanPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, notError(), cid, stx)
}

func getCanPulseGuide(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, true, cid, stx)
}

func getCanSetDeclinationRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSetGuideRates(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, true, cid, stx)
}

func getCanSetPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSetPierSide(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSetRightAscensionRate(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSetTracking(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, notError(), cid, stx)
}

func getCanSlew(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, true, cid, stx)
}

func getCanSlewAltAz(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSlewAltAzAsync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSlewAsync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, notError(), cid, stx)
}

func getCanSync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanSyncAltAz(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, false, cid, stx)
}

func getCanUnpark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeBool(c, notError(), cid, stx)
}

func getCanMoveAxis(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	axis, ok := paramInt(c, "Axis")
	if !ok {
		writeError(c, errInvalidValue, "Missing Axis", cid, stx)
		return
	}
	writeBool(c, (axis == 0 || axis == 1) && notError(), cid, stx)
}
