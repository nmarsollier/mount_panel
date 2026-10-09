package alpaca

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// response is the Alpaca response envelope. Value is raw JSON so numbers,
// booleans, strings and arrays are emitted exactly as built.
type response struct {
	Value               json.RawMessage `json:"Value,omitempty"`
	ClientTransactionID uint32          `json:"ClientTransactionID"`
	ServerTransactionID uint32          `json:"ServerTransactionID"`
	ErrorNumber         int             `json:"ErrorNumber"`
	ErrorMessage        string          `json:"ErrorMessage"`
}

var serverTx uint32

// nextServerTx returns a monotonically increasing server transaction id.
func nextServerTx() uint32 {
	return atomic.AddUint32(&serverTx, 1)
}

// writeValue sends a successful response with a raw JSON value.
func writeValue(c *gin.Context, raw json.RawMessage, cid, stx uint32) {
	c.JSON(http.StatusOK, response{
		Value:               raw,
		ClientTransactionID: cid,
		ServerTransactionID: stx,
	})
}

// writeOK sends a successful response with no value.
func writeOK(c *gin.Context, cid, stx uint32) {
	c.JSON(http.StatusOK, response{
		ClientTransactionID: cid,
		ServerTransactionID: stx,
	})
}

// writeError sends an Alpaca error response (always HTTP 200, like the spec).
func writeError(c *gin.Context, number int, message string, cid, stx uint32) {
	c.JSON(http.StatusOK, response{
		ClientTransactionID: cid,
		ServerTransactionID: stx,
		ErrorNumber:         number,
		ErrorMessage:        message,
	})
}

// writeBool sends a boolean value.
func writeBool(c *gin.Context, v bool, cid, stx uint32) {
	writeValue(c, rawBool(v), cid, stx)
}

// writeInt sends an integer value.
func writeInt(c *gin.Context, v int, cid, stx uint32) {
	writeValue(c, rawInt(v), cid, stx)
}

// writeFloat sends a float value formatted with six decimals (%.6f).
func writeFloat(c *gin.Context, v float64, cid, stx uint32) {
	writeValue(c, rawFloat(v), cid, stx)
}

// writeString sends a JSON string value.
func writeString(c *gin.Context, s string, cid, stx uint32) {
	writeValue(c, rawJSON(s), cid, stx)
}

// rawBool renders a boolean as a raw JSON value.
func rawBool(v bool) json.RawMessage {
	if v {
		return json.RawMessage("true")
	}
	return json.RawMessage("false")
}

// rawInt renders an integer as a raw JSON value.
func rawInt(v int) json.RawMessage {
	return json.RawMessage(strconv.Itoa(v))
}

// rawFloat renders a float with six decimals as a raw JSON value.
func rawFloat(v float64) json.RawMessage {
	return json.RawMessage(strconv.FormatFloat(v, 'f', 6, 64))
}

// rawJSON marshals any value to its JSON representation, ignoring errors.
func rawJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
