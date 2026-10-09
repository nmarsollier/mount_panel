package alpaca

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// putAbortSlew immediately stops any ongoing slew.
func putAbortSlew(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	if err := mount.Get().AbortSlew(); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putFindHome starts the home-finding sequence (asynchronous).
func putFindHome(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	if err := mount.Get().FindHome(); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putPark moves the mount to the park position.
func putPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	if err := mount.Get().Park(); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putUnpark takes the mount out of the parked state.
func putUnpark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	if err := mount.Get().Unpark(); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSetPark is not implemented — there is no configurable park position.
func putSetPark(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "SetPark is not supported", cid, stx)
}

// putSlewToCoordinatesAsync starts an asynchronous equatorial slew.
func putSlewToCoordinatesAsync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	ra, okRa := paramFloat(c, "RightAscension")
	dec, okDec := paramFloat(c, "Declination")
	if !okRa || !okDec {
		writeError(c, errInvalidValue, "Missing RightAscension or Declination", cid, stx)
		return
	}
	if err := mount.Get().SlewToCoordinates(ra, dec); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSlewToAltAzAsync is not implemented — equatorial mount.
func putSlewToAltAzAsync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Alt/Az slew not implemented (EQ mount)", cid, stx)
}

// putSyncToCoordinates is not supported by this mount.
func putSyncToCoordinates(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putSyncToAltAz is not supported by this mount.
func putSyncToAltAz(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putSlewToTargetAsync slews to the previously set target coordinates.
func putSlewToTargetAsync(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	if err := mount.Get().SlewToTarget(); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putSyncToTarget is not supported by this mount.
func putSyncToTarget(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()
	writeError(c, errNotImplemented, "Sync not supported by this mount", cid, stx)
}

// putPulseGuide executes a guide pulse. Direction is an integer (0-3).
func putPulseGuide(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	dir, okDir := paramInt(c, "Direction")
	dur, okDur := paramInt(c, "Duration")
	if !okDir || !okDur || dir < 0 || dir > 3 {
		writeError(c, errInvalidValue, "Missing or invalid Direction (0-3) or Duration", cid, stx)
		return
	}
	if dur <= 0 || dur > 30000 {
		writeError(c, errInvalidValue, "Duration out of range (1-30000 ms)", cid, stx)
		return
	}
	if err := mount.Get().PulseGuide(dir, dur); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}

// putMoveAxis moves a single axis continuously at the given rate in deg/s.
func putMoveAxis(c *gin.Context) {
	cid, stx := clientTransactionID(c), nextServerTx()

	axis, okAxis := paramInt(c, "Axis")
	rate, okRate := paramFloat(c, "Rate")
	if !okAxis || !okRate {
		writeError(c, errInvalidValue, "Missing Axis or Rate", cid, stx)
		return
	}
	if axis != 0 && axis != 1 {
		writeError(c, errInvalidValue, "Axis not supported", cid, stx)
		return
	}
	if err := mount.Get().MoveAxis(axis, rate); err != nil {
		writeError(c, errInvalidValue, err.Error(), cid, stx)
		return
	}
	writeOK(c, cid, stx)
}
