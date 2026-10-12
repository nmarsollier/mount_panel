package alpaca_mount

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// notError reports whether the mount is operational (not in the error state).
func notError() bool {
	isErr, err := mount.Get().IsError()
	if err != nil {
		// Unknown state: don't disable capabilities on a communication failure.
		return true
	}
	return !isErr
}

func getCanFindHome(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, notError(), cid, stx)
}

func getCanPark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, notError(), cid, stx)
}

func getCanPulseGuide(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, true, cid, stx)
}

func getCanSetDeclinationRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSetGuideRates(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, true, cid, stx)
}

func getCanSetPark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSetPierSide(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSetRightAscensionRate(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSetTracking(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, notError(), cid, stx)
}

func getCanSlew(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, true, cid, stx)
}

func getCanSlewAltAz(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSlewAltAzAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSlewAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, notError(), cid, stx)
}

func getCanSync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanSyncAltAz(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

func getCanUnpark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, notError(), cid, stx)
}

func getCanMoveAxis(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	axis, ok := server.ParamInt(c, "Axis")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing Axis", cid, stx)
		return
	}
	server.WriteBool(c, (axis == 0 || axis == 1) && notError(), cid, stx)
}
