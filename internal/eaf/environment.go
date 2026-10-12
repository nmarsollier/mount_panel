package eaf

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
	return env.HumidityPct, nil
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
	return env.DewPointC, nil
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
