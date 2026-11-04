package events

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// schema is the single append-only event table plus its one search index. The
// payload column holds the full event as compact JSON; the typed columns exist
// only to filter, order and correlate without parsing that JSON.
const schema = `
CREATE TABLE IF NOT EXISTS events (
    id             INTEGER PRIMARY KEY,
    ts             INTEGER NOT NULL,
    event_type     TEXT    NOT NULL,
    source         TEXT    NOT NULL,
    subject        TEXT,
    correlation_id TEXT,
    payload        TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_ts_type ON events (ts, event_type);
`

const insertSQL = `INSERT INTO events (ts, event_type, source, subject, correlation_id, payload) VALUES (?, ?, ?, ?, ?, ?)`

// Open opens (or creates) the event database at path and applies the schema.
// The connection is configured for WAL journaling and a busy timeout so an
// external SQLite client can read the file while the app writes. A single
// connection is used: the writer is the only internal writer, and one
// connection keeps the pragmas coherent and avoids self-contention.
func Open(path string) (*sql.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	for _, p := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
	} {
		if _, err := db.Exec(p); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
