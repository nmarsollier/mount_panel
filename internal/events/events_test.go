package events

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestRecordPersistsCorrelatedEvents exercises the full path: Init, two
// correlated Record calls, Close (flush) and a fresh Open that reads them back.
func TestRecordPersistsCorrelatedEvents(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.db")

	if err := Init(path); err != nil {
		t.Fatalf("Init: %v", err)
	}

	Record(Event{
		Type:          TypeAlpacaCommand,
		Source:        SourceAlpaca,
		Subject:       "/api/v1/telescope/0/slewtocoordinatesasync",
		CorrelationID: "c1",
		Payload:       map[string]any{"method": "PUT"},
	})
	Record(Event{
		Type:          TypeAlpacaResponse,
		Source:        SourceAlpaca,
		Subject:       "/api/v1/telescope/0/slewtocoordinatesasync",
		CorrelationID: "c1",
		Payload:       map[string]any{"status": 200},
	})

	Close()

	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	var n int
	if err := db.QueryRow(`SELECT count(*) FROM events`).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 2 {
		t.Fatalf("want 2 events, got %d", n)
	}

	var payload string
	if err := db.QueryRow(
		`SELECT payload FROM events WHERE event_type = ?`, TypeAlpacaCommand,
	).Scan(&payload); err != nil {
		t.Fatalf("select: %v", err)
	}
	if !json.Valid([]byte(payload)) {
		t.Fatalf("payload is not compact JSON: %s", payload)
	}
}
