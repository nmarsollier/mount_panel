package server

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nmarsollier/mount_panel/internal/events"
)

// Events records every Alpaca request/response pair as two correlated events:
// the command is recorded before the handler runs (before the internal services
// are called) and the response after it returns. Both share a correlation id so
// the pair can be joined, and the response body — always HTTP 200 with the
// error in the envelope — is captured for the audit trail.
func Events() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlation := uuid.NewString()

		events.Record(events.Event{
			Type:          events.TypeAlpacaCommand,
			Source:        events.SourceAlpaca,
			Subject:       c.Request.URL.Path,
			CorrelationID: correlation,
			Details: events.AlpacaCommand{
				Method: c.Request.Method,
				Path:   c.Request.URL.Path,
				Query:  c.Request.URL.RawQuery,
			},
		})

		bw := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = bw

		c.Next()

		var resp Response
		_ = json.Unmarshal(bw.body.Bytes(), &resp)

		events.Record(events.Event{
			Type:          events.TypeAlpacaResponse,
			Source:        events.SourceAlpaca,
			Subject:       c.Request.URL.Path,
			CorrelationID: correlation,
			Details: events.AlpacaResponse{
				Status:              c.Writer.Status(),
				ErrorNumber:         resp.ErrorNumber,
				ErrorMessage:        resp.ErrorMessage,
				ClientTransactionID: resp.ClientTransactionID,
				ServerTransactionID: resp.ServerTransactionID,
				Value:               resp.Value,
			},
		})
	}
}

// bodyWriter captures the response bytes while delegating everything else to
// the wrapped gin.ResponseWriter, so the Alpaca envelope can be audited.
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
