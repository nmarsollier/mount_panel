package eafproto

import "context"

// DeviceState is the general accessory snapshot: power, health and the
// subsystems summary.
type DeviceState struct {
	PowerExternalPresent bool    `json:"power_external_present"`
	Health               string  `json:"health"`
	FocuserBusy          bool    `json:"focuser_busy"`
	FanMode              string  `json:"fan_mode"`
	SensorHealth         string  `json:"sensor_health"`
	ResetReason          string  `json:"reset_reason"`
	Faults               []Fault `json:"faults"`
}

// Fault is a single reported fault, bounded in size on the accessory side.
type Fault struct {
	FaultID   string `json:"fault_id"`
	Subsystem string `json:"subsystem"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	FirstMs   int64  `json:"first_ms"`
	Active    bool   `json:"active"`
	Latched   bool   `json:"latched"`
}

// FocuserState is the focuser snapshot: the signed STEP counter, motion state
// and driver health.
type FocuserState struct {
	CounterSteps   int64   `json:"counter_steps"`
	CounterValid   bool    `json:"counter_valid"`
	IsMoving       bool    `json:"is_moving"`
	Source         string  `json:"source"`
	OpID           string  `json:"op_id"`
	OpStatus       string  `json:"op_status"`
	StartCounter   int64   `json:"start_counter"`
	TargetCounter  int64   `json:"target_counter"`
	ExecutedSteps  int64   `json:"executed_steps"`
	SpeedStepsS    float64 `json:"speed_steps_s"`
	EnabledCommand bool    `json:"enabled_command"`
	Guard          struct {
		Enabled    bool  `json:"enabled"`
		MinCounter int64 `json:"min_counter"`
		MaxCounter int64 `json:"max_counter"`
	} `json:"guard"`
	Driver struct {
		UartOK           bool  `json:"uart_ok"`
		LastReadMs       int64 `json:"last_read_ms"`
		OvertempWarning  bool  `json:"overtemp_warning"`
		OvertempShutdown bool  `json:"overtemp_shutdown"`
		ShortA           bool  `json:"short_a"`
		ShortB           bool  `json:"short_b"`
		OpenLoadA        bool  `json:"open_load_a"`
		OpenLoadB        bool  `json:"open_load_b"`
	} `json:"driver"`
}

// FanState is the fan snapshot: the selected mode and the outputs applied by
// the local thermal regulation.
type FanState struct {
	Mode           string   `json:"mode"`
	ModeRev        int      `json:"mode_rev"`
	Source         string   `json:"source"`
	OutputPWMPct   int      `json:"output_pwm_pct"`
	TargetPWMPct   int      `json:"target_pwm_pct"`
	RelayCommand   bool     `json:"relay_command"`
	PowerAvailable bool     `json:"power_available"`
	DeltaC         *float64 `json:"delta_c"`
	ThermalStage   string   `json:"thermal_stage"`
	Reason         string   `json:"reason"`
	LastChangeMs   int64    `json:"last_change_ms"`
}

// EnvironmentState is the environment snapshot: ambient/mirror temperatures,
// humidity and dew point (mock values until the sensors are ported).
type EnvironmentState struct {
	AmbientC            float64 `json:"ambient_c"`
	HumidityPct         float64 `json:"humidity_pct"`
	MirrorC             float64 `json:"mirror_c"`
	DewPointC           float64 `json:"dew_point_c"`
	MirrorAmbientDeltaC float64 `json:"mirror_ambient_delta_c"`
	MirrorDewMarginC    float64 `json:"mirror_dew_margin_c"`
}

// VibrationState is the vibration snapshot (ADXL345 statistics, m/s²). The
// firmware reports the last completed second's demeaned RMS.
type VibrationState struct {
	Enabled      bool    `json:"enabled"`
	Valid        bool    `json:"valid"`
	Error        *string `json:"error"`
	SampleRateHz int     `json:"sample_rate_hz"`
	RmsX         float64 `json:"rms_x"`
	RmsY         float64 `json:"rms_y"`
	RmsZ         float64 `json:"rms_z"`
	RmsTotal     float64 `json:"rms_total"`
	Processing   string  `json:"processing"`
}

// StateDevice asks the accessory for the general snapshot.
func (c *Client) StateDevice() (*DeviceState, error) {
	out := struct {
		Data DeviceState `json:"data"`
	}{}
	if err := c.request(context.Background(), "state", "device", "", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// StateFocuser asks the accessory for the focuser snapshot.
func (c *Client) StateFocuser() (*FocuserState, error) {
	out := struct {
		Data FocuserState `json:"data"`
	}{}
	if err := c.request(context.Background(), "state", "focuser", "", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// StateFan asks the accessory for the fan snapshot.
func (c *Client) StateFan() (*FanState, error) {
	out := struct {
		Data FanState `json:"data"`
	}{}
	if err := c.request(context.Background(), "state", "fan", "", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// StateEnvironment asks the accessory for the environment snapshot.
func (c *Client) StateEnvironment() (*EnvironmentState, error) {
	out := struct {
		Data EnvironmentState `json:"data"`
	}{}
	if err := c.request(context.Background(), "state", "environment", "", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// StateVibration asks the accessory for the vibration snapshot.
func (c *Client) StateVibration() (*VibrationState, error) {
	out := struct {
		Data VibrationState `json:"data"`
	}{}
	if err := c.request(context.Background(), "state", "vibration", "", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}
