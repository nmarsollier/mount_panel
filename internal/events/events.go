// Package events records every significant action in the gateway — Alpaca
// commands and their responses, serial commands sent to the mount, and sensor
// readings — into an embedded SQLite database as an append-only event log. Each
// event is an immutable JSON document plus a few typed columns (time, type,
// source, subject, correlation id) that make it searchable without parsing the
// JSON. Writes are asynchronous: Record only enqueues into a bounded channel
// that a single background writer drains in batches.
package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// Event type categories.
const (
	TypeAlpacaCommand  = "alpaca.command"
	TypeAlpacaResponse = "alpaca.response"
	TypeMountCommand   = "mount.command"
	TypeMountResponse  = "mount.response"
	TypeSensorReading  = "sensor.reading"
	TypeSystem         = "system"
)

// Event sources — the component that produced the event.
const (
	SourceAlpaca = "alpaca"
	SourceMount  = "mount"
	SourceEAF    = "eaf"
	SourceFan    = "fan"
	SourceAmb    = "amb"
	SourceUI     = "ui"
	SourceSystem = "system"
)

// Event is one immutable fact in the log. Payload carries the full event
// content and is marshaled to compact JSON (no indentation). Ts is the
// occurrence time in milliseconds since the Unix epoch (UTC); when zero it is
// set to now at Record time.
type Event struct {
	Ts            int64
	Type          string
	Source        string
	Subject       string
	CorrelationID string
	Payload       any
}

// row is the persisted shape of an event: payload is already a compact JSON
// string, so the writer never re-encodes on the hot path.
type row struct {
	ts            int64
	eventType     string
	source        string
	subject       string
	correlationID string
	payload       string
}

// Record queues an event for asynchronous persistence. It returns immediately;
// the disk write happens in the background writer. It is a no-op when Init has
// not been called (or Close has shut the store down), so the Alpaca and serial
// paths never fail because the audit log is unavailable.
func Record(e Event) {
	if e.Ts == 0 {
		e.Ts = time.Now().UnixMilli()
	}
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		logx.Warn(eventsCtx(), "event payload marshal failed: "+err.Error())
		payload = []byte("null")
	}
	record(row{
		ts:            e.Ts,
		eventType:     e.Type,
		source:        e.Source,
		subject:       e.Subject,
		correlationID: e.CorrelationID,
		payload:       string(payload),
	})
}

func eventsCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Events"))
}
