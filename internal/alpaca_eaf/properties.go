package alpaca_eaf

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/eaf"
	"github.com/nmarsollier/mount_panel/internal/fan"
	"github.com/nmarsollier/mount_panel/internal/temp"
)

// getAbsolute reports that the adapter publishes an absolute position.
func getAbsolute(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, true, cid, stx)
}

// getIsMoving reports whether the focuser is moving.
func getIsMoving(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := eaf.Get().IsMoving()
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v, cid, stx)
}

// getMaxIncrement returns the maximum per-move increment.
func getMaxIncrement(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, int(eaf.Get().MaxIncrement()), cid, stx)
}

// getMaxStep returns the maximum absolute position.
func getMaxStep(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, int(eaf.Get().MaxStep()), cid, stx)
}

// getPosition returns the current absolute position.
func getPosition(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := eaf.Get().Position()
	if err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteInt(c, int(v), cid, stx)
}

// getStepSize returns the calibrated step size, or NotImplemented while unknown.
func getStepSize(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "StepSize not calibrated", cid, stx)
}

// getTempComp reports whether temperature compensation is active (always false).
func getTempComp(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

// getTempCompAvailable reports whether temperature compensation is available.
func getTempCompAvailable(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

// getTemperature returns the ambient temperature in Celsius.
func getTemperature(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	v, err := temp.Get().AmbientTemperature()
	if err != nil {
		writeEnvironmentError(c, err, cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// writeEnvironmentError maps an environment read failure to Alpaca: an absent
// sensor is 1024 (not implemented), a broken serial link is 1027 (not
// connected).
func writeEnvironmentError(c *gin.Context, err error, cid, stx uint32) {
	if errors.Is(err, temp.ErrUnavailable) || errors.Is(err, fan.ErrUnavailable) {
		server.WriteError(c, server.ErrNotImplemented, err.Error(), cid, stx)
		return
	}
	server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
}
