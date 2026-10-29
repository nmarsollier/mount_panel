package fan

import "errors"

// ErrUnavailable reports that the requested reading has no valid value because
// its sensor is absent or faulted. Alpaca handlers map it to ErrorNumber 1024
// (not implemented) instead of a fabricated reading.
var ErrUnavailable = errors.New("fan reading unavailable")

// MirrorTemperature returns the primary-mirror temperature in Celsius.
func (f *Fan) MirrorTemperature() (float64, error) {
	c, err := f.getConn()
	if err != nil {
		return 0, err
	}
	env, err := c.StateEnvironment()
	if err != nil {
		return 0, err
	}
	if env.MirrorC == nil {
		return 0, ErrUnavailable
	}
	return *env.MirrorC, nil
}
