package eaf

import "errors"

// ErrEnvironmentUnavailable reports that the requested environment reading has
// no valid value because its sensor is absent or faulted. Alpaca handlers map
// it to ErrorNumber 1024 (not implemented) instead of a fabricated reading.
var ErrEnvironmentUnavailable = errors.New("environment reading unavailable")

// Humidity returns the ambient relative humidity in percent.
func (e *EAF) Humidity() (float64, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.HumidityPct == nil {
		return 0, ErrEnvironmentUnavailable
	}
	return *env.HumidityPct, nil
}

// DewPoint returns the dew point in Celsius.
func (e *EAF) DewPoint() (float64, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.DewPointC == nil {
		return 0, ErrEnvironmentUnavailable
	}
	return *env.DewPointC, nil
}

// MirrorTemperature returns the mirror temperature in Celsius.
func (e *EAF) MirrorTemperature() (float64, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.MirrorC == nil {
		return 0, ErrEnvironmentUnavailable
	}
	return *env.MirrorC, nil
}

// AveragePeriod returns the ObservingConditions averaging window, in hours.
func (e *EAF) AveragePeriod() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.averagePeriod
}

// SetAveragePeriod stores the ObservingConditions averaging window, in hours.
func (e *EAF) SetAveragePeriod(hours float64) {
	e.mu.Lock()
	e.averagePeriod = hours
	e.mu.Unlock()
}
