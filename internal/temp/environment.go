package temp

import "errors"

// ErrUnavailable reports that the requested reading has no valid value because
// its sensor is absent or faulted. Alpaca handlers map it to ErrorNumber 1024
// (not implemented) instead of a fabricated reading.
var ErrUnavailable = errors.New("environment reading unavailable")

// Humidity returns the ambient relative humidity in percent.
func (t *Temp) Humidity() (float64, error) {
	c, err := t.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.HumidityPct == nil {
		return 0, ErrUnavailable
	}
	return *env.HumidityPct, nil
}

// DewPoint returns the dew point in Celsius.
func (t *Temp) DewPoint() (float64, error) {
	c, err := t.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.DewPointC == nil {
		return 0, ErrUnavailable
	}
	return *env.DewPointC, nil
}

// AmbientTemperature returns the ambient temperature in Celsius.
func (t *Temp) AmbientTemperature() (float64, error) {
	c, err := t.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.AmbientC == nil {
		return 0, ErrUnavailable
	}
	return *env.AmbientC, nil
}
