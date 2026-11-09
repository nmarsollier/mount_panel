package events

import "encoding/json"

// This file defines the typed detail payloads, one per event type (and one per
// sensor source). Each maps to the columns of the matching detail table: the
// writer turns a detail into the ordered values of that table's INSERT. A
// pointer field means "this family did not carry it", and is written as NULL.

// AlpacaCommand is the detail of an alpaca.command event: the HTTP request the
// client sent, split into typed columns.
type AlpacaCommand struct {
	Method string
	Path   string
	Query  string
}

// AlpacaResponse is the detail of an alpaca.response event: the Alpaca Response
// envelope as columns. Value keeps the endpoint-specific result as raw JSON; it
// is the only non-scalar column of this table.
type AlpacaResponse struct {
	Status              int
	ErrorNumber         int
	ErrorMessage        string
	ClientTransactionID uint32
	ServerTransactionID uint32
	Value               json.RawMessage
}

// MountCommand is the detail of a mount.command event: the request envelope
// (family + message id) plus the union of request fields across the five mount
// message families. A pointer field is nil when the family did not carry it.
type MountCommand struct {
	Family string
	MsgID  uint32
	// config
	Lat          *float64
	Lon          *float64
	Elevation    *int
	GuideRateRa  *float64
	GuideRateDec *float64
	Utc          *string
	// control
	Tracking     string
	ManualRaDps  *float64
	ManualDecDps *float64
	// action
	Action     string
	Scope      string
	Ra         *float64
	Dec        *float64
	Speed      *int
	Direction  string
	DurationMs *uint32
	Axis       string
	Degrees    *float64
	Param      string
}

// MountResponse is the detail of a mount.response event: the response envelope
// (family, message id, ok, error) plus the union of response fields across the
// five families. A pointer field is nil for families that do not carry it.
type MountResponse struct {
	Family string
	MsgID  uint32
	OK     bool
	Error  string
	// capabilities
	BootID          *string
	Proto           *int
	DeviceType      *string
	DeviceID        *string
	Name            *string
	FirmwareVersion *string
	CapGuide        *bool
	CapHome         *bool
	CapPark         *bool
	CapUnpark       *bool
	CapMoveAxis     *bool
	CapSlew         *bool
	CapTracking     *bool
	CapSync         *bool
	CapMeridianFlip *bool
	CapPierSide     *bool
	Capabilities    json.RawMessage
	// state
	State       *string
	Tracking    *string
	RaSteps     *int64
	DecSteps    *int64
	RaAxisDeg   *float64
	DecAxisDeg  *float64
	Ra          *float64
	Dec         *float64
	Lst         *float64
	PierSide    *string
	RaSpeedDps  *float64
	DecSpeedDps *float64
	Guiding     *bool
	AtHome      *bool
	AtPark      *bool
	Power       *bool
	TimeValid   *bool
	RaMin       *float64
	RaMax       *float64
	DecMin      *float64
	DecMax      *float64
	ConfigRev   *uint32
	UptimeMs    *uint64
	// config
	Lat          *float64
	Lon          *float64
	Elevation    *int
	GuideRateRa  *float64
	GuideRateDec *float64
	Utc          *string
}

// SensorMount is the detail of a sensor.reading event from the mount, mirroring
// the mount status snapshot shown in the UI.
type SensorMount struct {
	Status            string
	Tracking          string
	Power             bool
	Ra                string
	Dec               string
	Lst               string
	PierSide          string
	Time              string
	IsHome            bool
	SettingsLat       float64
	SettingsLon       float64
	SettingsElevation int
	RaAxisDeg         float64
	DecAxisDeg        float64
	RaSteps           int64
	DecSteps          int64
	RaSpeed           float64
	DecSpeed          float64
	Guiding           bool
	Microsteps        int
	LimitsRaMin       float64
	LimitsRaMax       float64
	LimitsDecMin      float64
	LimitsDecMax      float64
	UptimeS           uint64
}

// SensorEAF is the detail of a sensor.reading event from the focuser.
type SensorEAF struct {
	Detected     bool
	BootID       string
	Power        bool
	Health       string
	CounterSteps int64
	IsMoving     bool
	OpStatus     string
	DriverOK     bool
	Enabled      bool
}

// SensorFan is the detail of a sensor.reading event from the fan controller.
type SensorFan struct {
	Detected       bool
	BootID         string
	Power          bool
	Health         string
	FanError       *string
	FanMode        string
	FanPower       bool
	FanSpeed       int
	FanManualSpeed int
	FanPowerOK     bool
	FanRPM         *int
	MirrorError    *string
	MirrorC        *float64
}

// SensorAmb is the detail of a sensor.reading event from the ambient controller.
type SensorAmb struct {
	Detected       bool
	BootID         string
	Health         string
	SensorHealth   string
	EnvError       *string
	AmbientC       *float64
	HumidityPct    *float64
	DewPointC      *float64
	VibError       *string
	RmsTotal       *float64
	AccelPeakTotal *float64
	VelRmsTotal    *float64
	VelPeakTotal   *float64
	IsoZone        *string
	OriError       *string
	PitchDeg       *float64
	RollDeg        *float64
	AxG            *float64
	AyG            *float64
	AzG            *float64
}
