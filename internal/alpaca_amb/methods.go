package alpaca_amb

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
)

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
