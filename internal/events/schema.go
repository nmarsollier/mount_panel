package events

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// schemaVersion is bumped whenever the schema changes shape. On Open, a file
// whose user_version differs is treated as obsolete: its event tables are
// dropped and rebuilt from scratch. The log is an audit trail, not durable
// state worth migrating.
const schemaVersion = 1

// eventTables lists every table owned by this package, in an order that drops
// children before the parent (matters only if foreign keys are ever enforced).
var eventTables = []string{
	"event_sensor_amb",
	"event_sensor_fan",
	"event_sensor_eaf",
	"event_sensor_mount",
	"event_mount_response",
	"event_mount_command",
	"event_alpaca_response",
	"event_alpaca_command",
	"events",
}

// schema is the shared events table (header facts only) plus one typed detail
// table per event type. Each detail row maps to its parent one-to-(zero-or-one)
// through event_id; the typed columns hold the event data instead of JSON.
const schema = `
CREATE TABLE IF NOT EXISTS events (
    id             INTEGER PRIMARY KEY,
    ts             INTEGER NOT NULL,
    event_type     TEXT    NOT NULL,
    source         TEXT    NOT NULL,
    subject        TEXT,
    correlation_id TEXT
);
CREATE INDEX IF NOT EXISTS idx_events_ts_type ON events (ts, event_type);

CREATE TABLE IF NOT EXISTS event_alpaca_command (
    event_id INTEGER PRIMARY KEY REFERENCES events(id),
    method   TEXT NOT NULL,
    path     TEXT NOT NULL,
    query    TEXT
);

CREATE TABLE IF NOT EXISTS event_alpaca_response (
    event_id              INTEGER PRIMARY KEY REFERENCES events(id),
    status                INTEGER NOT NULL,
    error_number          INTEGER NOT NULL,
    error_message         TEXT,
    client_transaction_id INTEGER,
    server_transaction_id INTEGER,
    value                 TEXT
);

CREATE TABLE IF NOT EXISTS event_mount_command (
    event_id       INTEGER PRIMARY KEY REFERENCES events(id),
    family         TEXT    NOT NULL,
    msg_id         INTEGER NOT NULL,
    lat            REAL,
    lon            REAL,
    elevation      INTEGER,
    guide_rate_ra  REAL,
    guide_rate_dec REAL,
    utc            TEXT,
    tracking       TEXT,
    manual_ra_dps  REAL,
    manual_dec_dps REAL,
    action         TEXT,
    scope          TEXT,
    ra             REAL,
    dec            REAL,
    speed          INTEGER,
    direction      TEXT,
    duration_ms    INTEGER,
    axis           TEXT,
    degrees        REAL,
    param          TEXT
);

CREATE TABLE IF NOT EXISTS event_mount_response (
    event_id          INTEGER PRIMARY KEY REFERENCES events(id),
    family            TEXT    NOT NULL,
    msg_id            INTEGER NOT NULL,
    ok                INTEGER NOT NULL,
    error             TEXT,
    boot_id           TEXT,
    proto             INTEGER,
    device_type       TEXT,
    device_id         TEXT,
    name              TEXT,
    firmware_version  TEXT,
    cap_guide         INTEGER,
    cap_home          INTEGER,
    cap_park          INTEGER,
    cap_unpark        INTEGER,
    cap_move_axis     INTEGER,
    cap_slew          INTEGER,
    cap_tracking      INTEGER,
    cap_sync          INTEGER,
    cap_meridian_flip INTEGER,
    cap_pier_side     INTEGER,
    capabilities_json TEXT,
    state             TEXT,
    tracking          TEXT,
    ra_steps          INTEGER,
    dec_steps         INTEGER,
    ra_axis_deg       REAL,
    dec_axis_deg      REAL,
    ra                REAL,
    dec               REAL,
    lst               REAL,
    pier_side         TEXT,
    ra_speed_dps      REAL,
    dec_speed_dps     REAL,
    guiding           INTEGER,
    at_home           INTEGER,
    at_park           INTEGER,
    power             INTEGER,
    time_valid        INTEGER,
    ra_min            REAL,
    ra_max            REAL,
    dec_min           REAL,
    dec_max           REAL,
    config_rev        INTEGER,
    uptime_ms         INTEGER,
    lat               REAL,
    lon               REAL,
    elevation         INTEGER,
    guide_rate_ra     REAL,
    guide_rate_dec    REAL,
    utc               TEXT
);

CREATE TABLE IF NOT EXISTS event_sensor_mount (
    event_id           INTEGER PRIMARY KEY REFERENCES events(id),
    status             TEXT,
    tracking           TEXT,
    power              INTEGER,
    ra                 TEXT,
    dec                TEXT,
    lst                TEXT,
    pier_side          TEXT,
    time               TEXT,
    is_home            INTEGER,
    settings_lat       REAL,
    settings_lon       REAL,
    settings_elevation INTEGER,
    ra_axis_deg        REAL,
    dec_axis_deg       REAL,
    ra_steps           INTEGER,
    dec_steps          INTEGER,
    ra_speed           REAL,
    dec_speed          REAL,
    guiding            INTEGER,
    microsteps         INTEGER,
    limits_ra_min      REAL,
    limits_ra_max      REAL,
    limits_dec_min     REAL,
    limits_dec_max     REAL,
    uptime_s           INTEGER
);

CREATE TABLE IF NOT EXISTS event_sensor_eaf (
    event_id      INTEGER PRIMARY KEY REFERENCES events(id),
    detected      INTEGER,
    boot_id       TEXT,
    power         INTEGER,
    health        TEXT,
    counter_steps INTEGER,
    is_moving     INTEGER,
    op_status     TEXT,
    driver_ok     INTEGER,
    enabled       INTEGER
);

CREATE TABLE IF NOT EXISTS event_sensor_fan (
    event_id         INTEGER PRIMARY KEY REFERENCES events(id),
    detected         INTEGER,
    boot_id          TEXT,
    power            INTEGER,
    health           TEXT,
    fan_error        TEXT,
    fan_mode         TEXT,
    fan_power        INTEGER,
    fan_speed        INTEGER,
    fan_manual_speed INTEGER,
    fan_power_ok     INTEGER,
    fan_rpm          INTEGER,
    mirror_error     TEXT,
    mirror_c         REAL
);

CREATE TABLE IF NOT EXISTS event_sensor_amb (
    event_id         INTEGER PRIMARY KEY REFERENCES events(id),
    detected         INTEGER,
    boot_id          TEXT,
    health           TEXT,
    sensor_health    TEXT,
    env_error        TEXT,
    ambient_c        REAL,
    humidity_pct     REAL,
    dew_point_c      REAL,
    vib_error        TEXT,
    rms_total        REAL,
    accel_peak_total REAL,
    vel_rms_total    REAL,
    vel_peak_total   REAL,
    iso_zone         TEXT,
    ori_error        TEXT,
    pitch_deg        REAL,
    roll_deg         REAL,
    ax_g             REAL,
    ay_g             REAL,
    az_g             REAL
);
`

const insertEventSQL = `INSERT INTO events (ts, event_type, source, subject, correlation_id) VALUES (?, ?, ?, ?, ?)`

const insertAlpacaCommandSQL = `INSERT INTO event_alpaca_command (event_id, method, path, query) VALUES (?, ?, ?, ?)`

const insertAlpacaResponseSQL = `INSERT INTO event_alpaca_response (event_id, status, error_number, error_message, client_transaction_id, server_transaction_id, value) VALUES (?, ?, ?, ?, ?, ?, ?)`

const insertMountCommandSQL = `INSERT INTO event_mount_command (event_id, family, msg_id, lat, lon, elevation, guide_rate_ra, guide_rate_dec, utc, tracking, manual_ra_dps, manual_dec_dps, action, scope, ra, dec, speed, direction, duration_ms, axis, degrees, param) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertMountResponseSQL = `INSERT INTO event_mount_response (event_id, family, msg_id, ok, error, boot_id, proto, device_type, device_id, name, firmware_version, cap_guide, cap_home, cap_park, cap_unpark, cap_move_axis, cap_slew, cap_tracking, cap_sync, cap_meridian_flip, cap_pier_side, capabilities_json, state, tracking, ra_steps, dec_steps, ra_axis_deg, dec_axis_deg, ra, dec, lst, pier_side, ra_speed_dps, dec_speed_dps, guiding, at_home, at_park, power, time_valid, ra_min, ra_max, dec_min, dec_max, config_rev, uptime_ms, lat, lon, elevation, guide_rate_ra, guide_rate_dec, utc) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertSensorMountSQL = `INSERT INTO event_sensor_mount (event_id, status, tracking, power, ra, dec, lst, pier_side, time, is_home, settings_lat, settings_lon, settings_elevation, ra_axis_deg, dec_axis_deg, ra_steps, dec_steps, ra_speed, dec_speed, guiding, microsteps, limits_ra_min, limits_ra_max, limits_dec_min, limits_dec_max, uptime_s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertSensorEafSQL = `INSERT INTO event_sensor_eaf (event_id, detected, boot_id, power, health, counter_steps, is_moving, op_status, driver_ok, enabled) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertSensorFanSQL = `INSERT INTO event_sensor_fan (event_id, detected, boot_id, power, health, fan_error, fan_mode, fan_power, fan_speed, fan_manual_speed, fan_power_ok, fan_rpm, mirror_error, mirror_c) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertSensorAmbSQL = `INSERT INTO event_sensor_amb (event_id, detected, boot_id, health, sensor_health, env_error, ambient_c, humidity_pct, dew_point_c, vib_error, rms_total, accel_peak_total, vel_rms_total, vel_peak_total, iso_zone, ori_error, pitch_deg, roll_deg, ax_g, ay_g, az_g) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

	if err := applySchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// applySchema reconciles the on-disk schema with the current one. When the
// file's user_version does not match, the event tables are dropped and rebuilt
// from scratch before the current schema is applied.
func applySchema(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}

	if version != schemaVersion {
		for _, table := range eventTables {
			if _, err := db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
				return err
			}
		}
	}

	if _, err := db.Exec(schema); err != nil {
		return err
	}
	_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", schemaVersion))
	return err
}
