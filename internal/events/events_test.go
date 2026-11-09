package events

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestRecordPersistsTypedDetails exercises the full path: Init, several Record
// calls with typed details, Close (flush) and a fresh Open that reads them back
// from the shared table and the typed detail tables.
func TestRecordPersistsTypedDetails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.db")

	if err := Init(path); err != nil {
		t.Fatalf("Init: %v", err)
	}

	Record(Event{
		Type:          TypeAlpacaCommand,
		Source:        SourceAlpaca,
		Subject:       "/api/v1/telescope/0/slewtocoordinatesasync",
		CorrelationID: "c1",
		Details: AlpacaCommand{
			Method: "PUT",
			Path:   "/api/v1/telescope/0/slewtocoordinatesasync",
		},
	})
	Record(Event{
		Type:          TypeAlpacaResponse,
		Source:        SourceAlpaca,
		Subject:       "/api/v1/telescope/0/slewtocoordinatesasync",
		CorrelationID: "c1",
		Details: AlpacaResponse{
			Status:      200,
			ErrorNumber: 0,
			Value:       json.RawMessage(`true`),
		},
	})
	Record(Event{
		Type:    TypeSensorReading,
		Source:  SourceMount,
		Subject: "status",
		Details: SensorMount{Status: "tracking", Ra: "12:34:56"},
	})

	Close()

	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	var n int
	if err := db.QueryRow(`SELECT count(*) FROM events`).Scan(&n); err != nil {
		t.Fatalf("count events: %v", err)
	}
	if n != 3 {
		t.Fatalf("want 3 events, got %d", n)
	}

	var method string
	if err := db.QueryRow(
		`SELECT method FROM event_alpaca_command WHERE event_id = 1`,
	).Scan(&method); err != nil {
		t.Fatalf("select command: %v", err)
	}
	if method != "PUT" {
		t.Fatalf("want method PUT, got %s", method)
	}

	var status int
	if err := db.QueryRow(
		`SELECT status FROM event_alpaca_response WHERE event_id = 2`,
	).Scan(&status); err != nil {
		t.Fatalf("select response: %v", err)
	}
	if status != 200 {
		t.Fatalf("want status 200, got %d", status)
	}

	var ra string
	if err := db.QueryRow(
		`SELECT ra FROM event_sensor_mount WHERE event_id = 3`,
	).Scan(&ra); err != nil {
		t.Fatalf("select sensor: %v", err)
	}
	if ra != "12:34:56" {
		t.Fatalf("want ra 12:34:56, got %s", ra)
	}
}
