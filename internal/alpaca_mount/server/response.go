package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// Response is the Alpaca Response envelope. Value is raw JSON so numbers,
// booleans, strings and arrays are emitted exactly as built.
type Response struct {
	Value               json.RawMessage `json:"Value,omitempty"`
	ClientTransactionID uint32          `json:"ClientTransactionID"`
	ServerTransactionID uint32          `json:"ServerTransactionID"`
	ErrorNumber         int             `json:"ErrorNumber"`
	ErrorMessage        string          `json:"ErrorMessage"`
}

var serverTx uint32

// NextServerTx returns a monotonically increasing server transaction id.
func NextServerTx() uint32 {
	return atomic.AddUint32(&serverTx, 1)
}

// WriteValue sends a successful Response with a raw JSON value.
func WriteValue(c *gin.Context, raw json.RawMessage, cid, stx uint32) {
	c.JSON(http.StatusOK, Response{
		Value:               raw,
		ClientTransactionID: cid,
		ServerTransactionID: stx,
	})
}

// WriteOK sends a successful Response with no value.
func WriteOK(c *gin.Context, cid, stx uint32) {
	c.JSON(http.StatusOK, Response{
		ClientTransactionID: cid,
		ServerTransactionID: stx,
	})
}

// WriteError sends an Alpaca error Response (always HTTP 200, like the spec).
func WriteError(c *gin.Context, number int, message string, cid, stx uint32) {
	c.JSON(http.StatusOK, Response{
		ClientTransactionID: cid,
		ServerTransactionID: stx,
		ErrorNumber:         number,
		ErrorMessage:        message,
	})
}

// WriteBool sends a boolean value.
func WriteBool(c *gin.Context, v bool, cid, stx uint32) {
	WriteValue(c, RawBool(v), cid, stx)
}

// WriteInt sends an integer value.
func WriteInt(c *gin.Context, v int, cid, stx uint32) {
	WriteValue(c, RawInt(v), cid, stx)
}

// WriteFloat sends a float value formatted with six decimals (%.6f).
func WriteFloat(c *gin.Context, v float64, cid, stx uint32) {
	WriteValue(c, RawFloat(v), cid, stx)
}

// WriteString sends a JSON string value.
func WriteString(c *gin.Context, s string, cid, stx uint32) {
	WriteValue(c, RawJSON(s), cid, stx)
}

// RawBool renders a boolean as a raw JSON value.
func RawBool(v bool) json.RawMessage {
	if v {
		return json.RawMessage("true")
	}
	return json.RawMessage("false")
}

// RawInt renders an integer as a raw JSON value.
func RawInt(v int) json.RawMessage {
	return json.RawMessage(strconv.Itoa(v))
}

// RawFloat renders a float with six decimals as a raw JSON value.
func RawFloat(v float64) json.RawMessage {
	return json.RawMessage(strconv.FormatFloat(v, 'f', 6, 64))
}

// RawJSON marshals any value to its JSON representation, ignoring errors.
func RawJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
