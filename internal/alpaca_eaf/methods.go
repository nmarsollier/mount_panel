package alpaca_eaf

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"github.com/nmarsollier/nrm_ha_go/internal/eaf"
)

// putMove moves the focuser to an absolute position.
func putMove(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()

	position, ok := server.ParamInt(c, "Position")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing Position", cid, stx)
		return
	}
	if err := eaf.Get().MoveTo(int32(position)); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putHalt stops the focuser.
func putHalt(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	if err := eaf.Get().Halt(); err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// putTempComp is not implemented.
func putTempComp(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "TempComp not implemented", cid, stx)
}

// putCommandBlind is a legacy method, not supported.
func putCommandBlind(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "CommandBlind not supported", cid, stx)
}

// putCommandBool is a legacy method, not supported.
func putCommandBool(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "CommandBool not supported", cid, stx)
}

// putCommandString is a legacy method, not supported.
func putCommandString(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "CommandString not supported", cid, stx)
}

// putAction is not supported (SupportedActions is empty).
func putAction(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Action not supported", cid, stx)
}
