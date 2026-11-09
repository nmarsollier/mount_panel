// Package events records every significant action in the gateway — Alpaca
// commands and their responses, serial commands sent to the mount, and sensor
// readings — into an embedded SQLite database as an append-only event log. Each
// event is one row in the shared `events` table (time, type, source, subject,
// correlation id) plus, for most types, a typed row in the type-specific detail
// table that maps one-to-(zero-or-one) back to it. The detail tables hold the
// event's data as scalar columns instead of a JSON blob, so an external SQLite
// client can read and filter the log without parsing JSON. Writes are
// asynchronous: Record only enqueues into a bounded channel that a single
// background writer drains in batches.
package events

import (
	"context"
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

// Event is one immutable fact in the log. The header fields go into the shared
// events table; Details is the typed payload for the type-specific detail table
// (see details.go). Ts is the occurrence time in milliseconds since the Unix
// epoch (UTC); when zero it is set to now at Record time.
type Event struct {
	Ts            int64
	Type          string
	Source        string
	Subject       string
	CorrelationID string
	Details       any
}

// row is the persisted shape of an event awaiting a disk write.
type row struct {
	ts            int64
	eventType     string
	source        string
	subject       string
	correlationID string
	details       any
}

// Record queues an event for asynchronous persistence. It returns immediately;
// the disk write happens in the background writer. It is a no-op when Init has
// not been called (or Close has shut the store down), so the Alpaca and serial
// paths never fail because the audit log is unavailable.
func Record(e Event) {
	if e.Ts == 0 {
		e.Ts = time.Now().UnixMilli()
	}
	record(row{
		ts:            e.Ts,
		eventType:     e.Type,
		source:        e.Source,
		subject:       e.Subject,
		correlationID: e.CorrelationID,
		details:       e.Details,
	})
}

func eventsCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Events"))
}
