package events

import (
	"database/sql"
	"sync"
	"time"

	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

const (
	channelCapacity = 8192
	batchSize       = 256
	flushInterval   = 100 * time.Millisecond
)

var (
	mu     sync.RWMutex // guards db, ch and closed
	db     *sql.DB
	ch     chan row
	closed bool

	wg sync.WaitGroup
)

// Init opens the event database at path, applies the schema and starts the
// background writer. It is idempotent; a second call is a no-op.
func Init(path string) error {
	mu.Lock()
	defer mu.Unlock()

	if ch != nil {
		return nil
	}

	d, err := Open(path)
	if err != nil {
		return err
	}
	db = d
	closed = false
	ch = make(chan row, channelCapacity)
	wg.Add(1)
	go writer(ch, d)
	return nil
}

// Close flushes pending events, stops the writer and closes the database.
func Close() {
	mu.Lock()
	if closed {
		mu.Unlock()
		return
	}
	closed = true
	c := ch
	mu.Unlock()

	if c == nil {
		return
	}
	close(c)
	wg.Wait()

	mu.Lock()
	_ = db.Close()
	db = nil
	ch = nil
	mu.Unlock()
}

// record enqueues one event. It holds the read lock while sending so that Close,
// which needs the write lock to close the channel, waits for every in-flight
// send to finish and can never close a channel still being written to.
func record(r row) {
	mu.RLock()
	defer mu.RUnlock()

	if closed || ch == nil {
		return
	}
	ch <- r
}

// writer drains the channel and persists events in batches, one transaction per
// batch, so the amortized cost per event stays low under sustained load.
func writer(ch <-chan row, db *sql.DB) {
	defer wg.Done()

	var batch []row
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := insert(db, batch); err != nil {
			logx.Warn(eventsCtx(), "event insert failed: "+err.Error())
		}
		batch = batch[:0]
	}

	for {
		select {
		case r, ok := <-ch:
			if !ok {
				flush()
				return
			}
			batch = append(batch, r)
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

// insert writes a batch of rows inside one transaction.
func insert(db *sql.DB, rows []row) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, r := range rows {
		if _, err := stmt.Exec(r.ts, r.eventType, r.source, r.subject, r.correlationID, r.payload); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
