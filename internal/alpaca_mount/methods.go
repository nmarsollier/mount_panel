package alpaca_mount

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/mount"
)

// putAbortSlew immediately stops any ongoing slew.
func putAbortSlew(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := mount.Get().AbortSlew(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putFindHome starts the home-finding sequence (asynchronous).
func putFindHome(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := mount.Get().FindHome(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putPark moves the mount to the park position.
func putPark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := mount.Get().Park(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putUnpark takes the mount out of the parked state.
func putUnpark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := mount.Get().Unpark(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSetPark is not implemented — there is no configurable park position.
func putSetPark(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "SetPark is not supported", cid, stx)
}

// putSlewToCoordinatesAsync starts an asynchronous equatorial slew.
func putSlewToCoordinatesAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	ra, okRa := server.ParamFloat(c, "RightAscension")
	dec, okDec := server.ParamFloat(c, "Declination")
	if !okRa || !okDec {
		server.WriteError(c, server.ErrInvalidValue, "Missing RightAscension or Declination", cid, stx)
		return
	}
	if err := mount.Get().SlewToCoordinates(ra, dec); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSlewToAltAzAsync is not implemented — equatorial mount.
func putSlewToAltAzAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Alt/Az slew not implemented (EQ mount)", cid, stx)
}

// putSyncToCoordinates is not supported by this mount.
func putSyncToCoordinates(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putSyncToAltAz is not supported by this mount.
func putSyncToAltAz(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putSlewToTargetAsync slews to the previously set target coordinates.
func putSlewToTargetAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := mount.Get().SlewToTarget(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putSyncToTarget is not supported by this mount.
func putSyncToTarget(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putPulseGuide executes a guide pulse. Direction is an integer (0-3).
func putPulseGuide(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	dir, okDir := server.ParamInt(c, "Direction")
	dur, okDur := server.ParamInt(c, "Duration")
	if !okDir || !okDur || dir < 0 || dir > 3 {
		server.WriteError(c, server.ErrInvalidValue, "Missing or invalid Direction (0-3) or Duration", cid, stx)
		return
	}
	if dur <= 0 || dur > 30000 {
		server.WriteError(c, server.ErrInvalidValue, "Duration out of range (1-30000 ms)", cid, stx)
		return
	}
	if err := mount.Get().PulseGuide(dir, dur); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putMoveAxis moves a single axis continuously at the given rate in deg/s.
func putMoveAxis(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	axis, okAxis := server.ParamInt(c, "Axis")
	rate, okRate := server.ParamFloat(c, "Rate")
	if !okAxis || !okRate {
		server.WriteError(c, server.ErrInvalidValue, "Missing Axis or Rate", cid, stx)
		return
	}
	if axis != 0 && axis != 1 {
		server.WriteError(c, server.ErrInvalidValue, "Axis not supported", cid, stx)
		return
	}
	if err := mount.Get().MoveAxis(axis, rate); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}
