package alpaca

import (
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseForm ensures the request form (query string + URL-encoded body) is
// available. Parameter lookups then read query string first, body second,
// matching the firmware's precedence.
func parseForm(c *gin.Context) {
	_ = c.Request.ParseForm()
}

// param returns a parameter from the query string, falling back to the form
// body. Returns "" when absent.
func param(c *gin.Context, key string) string {
	parseForm(c)
	return c.Request.Form.Get(key)
}

// paramRaw returns a parameter string and whether it was present.
func paramRaw(c *gin.Context, key string) (string, bool) {
	v := param(c, key)
	return v, v != ""
}

// clientTransactionID extracts and strictly parses the ClientTransactionID,
// returning 0 when absent or malformed.
func clientTransactionID(c *gin.Context) uint32 {
	v := param(c, "ClientTransactionID")
	if v == "" {
		return 0
	}
	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(n)
}

// paramFloat parses a finite float parameter.
func paramFloat(c *gin.Context, key string) (float64, bool) {
	v := param(c, key)
	if v == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, false
	}
	return f, true
}

// paramInt parses an integer parameter.
func paramInt(c *gin.Context, key string) (int, bool) {
	v := param(c, key)
	if v == "" {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// paramBool parses a boolean parameter ("true"/"false", case-insensitive, or
// "1"/"0"). A malformed value is rejected, not treated as false.
func paramBool(c *gin.Context, key string) (bool, bool) {
	switch strings.ToLower(param(c, key)) {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	default:
		return false, false
	}
}
