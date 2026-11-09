package events

import (
	"database/sql"
	"encoding/json"
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

// statements holds the prepared insert statements, one per table. The parent
// statement writes the shared events row; the rest write the typed detail rows.
type statements struct {
	parent         *sql.Stmt
	alpacaCommand  *sql.Stmt
	alpacaResponse *sql.Stmt
	mountCommand   *sql.Stmt
	mountResponse  *sql.Stmt
	sensorMount    *sql.Stmt
	sensorEAF      *sql.Stmt
	sensorFan      *sql.Stmt
	sensorAmb      *sql.Stmt
}

// prepareAll prepares every insert statement on the transaction.
func prepareAll(tx *sql.Tx) (*statements, error) {
	var s statements
	targets := []struct {
		dst **sql.Stmt
		sql string
	}{
		{&s.parent, insertEventSQL},
		{&s.alpacaCommand, insertAlpacaCommandSQL},
		{&s.alpacaResponse, insertAlpacaResponseSQL},
		{&s.mountCommand, insertMountCommandSQL},
		{&s.mountResponse, insertMountResponseSQL},
		{&s.sensorMount, insertSensorMountSQL},
		{&s.sensorEAF, insertSensorEafSQL},
		{&s.sensorFan, insertSensorFanSQL},
		{&s.sensorAmb, insertSensorAmbSQL},
	}
	for _, t := range targets {
		st, err := tx.Prepare(t.sql)
		if err != nil {
			return nil, err
		}
		*t.dst = st
	}
	return &s, nil
}

func (s *statements) close() {
	for _, st := range []*sql.Stmt{
		s.parent, s.alpacaCommand, s.alpacaResponse, s.mountCommand,
		s.mountResponse, s.sensorMount, s.sensorEAF, s.sensorFan, s.sensorAmb,
	} {
		if st != nil {
			_ = st.Close()
		}
	}
}

// insert writes a batch of events inside one transaction: the shared row first,
// then its typed detail row, keyed by the parent's generated id.
func insert(db *sql.DB, rows []row) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmts, err := prepareAll(tx)
	if err != nil {
		return err
	}
	defer stmts.close()

	for _, r := range rows {
		res, err := stmts.parent.Exec(r.ts, r.eventType, r.source, r.subject, r.correlationID)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		child, values := childInsert(stmts, r)
		if child == nil {
			continue
		}
		args := append([]any{id}, values...)
		if _, err := child.Exec(args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// childInsert returns the prepared statement and ordered column values for the
// row's typed detail, or nil when the type has no detail (or none was set).
func childInsert(s *statements, r row) (*sql.Stmt, []any) {
	if r.details == nil {
		return nil, nil
	}

	switch r.eventType {
	case TypeAlpacaCommand:
		d := r.details.(AlpacaCommand)
		return s.alpacaCommand, []any{d.Method, d.Path, d.Query}
	case TypeAlpacaResponse:
		d := r.details.(AlpacaResponse)
		return s.alpacaResponse, []any{d.Status, d.ErrorNumber, d.ErrorMessage, d.ClientTransactionID, d.ServerTransactionID, rawOrNil(d.Value)}
	case TypeMountCommand:
		d := r.details.(MountCommand)
		return s.mountCommand, []any{
			d.Family, d.MsgID,
			d.Lat, d.Lon, d.Elevation, d.GuideRateRa, d.GuideRateDec, d.Utc,
			d.Tracking, d.ManualRaDps, d.ManualDecDps,
			d.Action, d.Scope, d.Ra, d.Dec, d.Speed, d.Direction, d.DurationMs, d.Axis, d.Degrees, d.Param,
		}
	case TypeMountResponse:
		d := r.details.(MountResponse)
		return s.mountResponse, []any{
			d.Family, d.MsgID, d.OK, d.Error,
			d.BootID, d.Proto, d.DeviceType, d.DeviceID, d.Name, d.FirmwareVersion,
			d.CapGuide, d.CapHome, d.CapPark, d.CapUnpark, d.CapMoveAxis, d.CapSlew,
			d.CapTracking, d.CapSync, d.CapMeridianFlip, d.CapPierSide, rawOrNil(d.Capabilities),
			d.State, d.Tracking, d.RaSteps, d.DecSteps, d.RaAxisDeg, d.DecAxisDeg,
			d.Ra, d.Dec, d.Lst, d.PierSide, d.RaSpeedDps, d.DecSpeedDps,
			d.Guiding, d.AtHome, d.AtPark, d.Power, d.TimeValid,
			d.RaMin, d.RaMax, d.DecMin, d.DecMax, d.ConfigRev, d.UptimeMs,
			d.Lat, d.Lon, d.Elevation, d.GuideRateRa, d.GuideRateDec, d.Utc,
		}
	case TypeSensorReading:
		switch r.source {
		case SourceMount:
			d := r.details.(SensorMount)
			return s.sensorMount, []any{
				d.Status, d.Tracking, d.Power, d.Ra, d.Dec, d.Lst, d.PierSide, d.Time, d.IsHome,
				d.SettingsLat, d.SettingsLon, d.SettingsElevation,
				d.RaAxisDeg, d.DecAxisDeg, d.RaSteps, d.DecSteps, d.RaSpeed, d.DecSpeed,
				d.Guiding, d.Microsteps,
				d.LimitsRaMin, d.LimitsRaMax, d.LimitsDecMin, d.LimitsDecMax, d.UptimeS,
			}
		case SourceEAF:
			d := r.details.(SensorEAF)
			return s.sensorEAF, []any{
				d.Detected, d.BootID, d.Power, d.Health,
				d.CounterSteps, d.IsMoving, d.OpStatus, d.DriverOK, d.Enabled,
			}
		case SourceFan:
			d := r.details.(SensorFan)
			return s.sensorFan, []any{
				d.Detected, d.BootID, d.Power, d.Health,
				d.FanError, d.FanMode, d.FanPower, d.FanSpeed, d.FanManualSpeed, d.FanPowerOK, d.FanRPM,
				d.MirrorError, d.MirrorC,
			}
		case SourceAmb:
			d := r.details.(SensorAmb)
			return s.sensorAmb, []any{
				d.Detected, d.BootID, d.Health, d.SensorHealth,
				d.EnvError, d.AmbientC, d.HumidityPct, d.DewPointC,
				d.VibError, d.RmsTotal, d.AccelPeakTotal, d.VelRmsTotal, d.VelPeakTotal, d.IsoZone,
				d.OriError, d.PitchDeg, d.RollDeg, d.AxG, d.AyG, d.AzG,
			}
		}
	}
	return nil, nil
}

// rawOrNil renders a raw JSON value as a bound []byte, or NULL when empty.
func rawOrNil(v json.RawMessage) any {
	if len(v) == 0 {
		return nil
	}
	return []byte(v)
}
