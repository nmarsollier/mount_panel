package eaf

// FanMode returns the current fan mode ("off", "auto" or "manual").
func (e *EAF) FanMode() (string, error) {
	c, err := e.getConn()
	if err != nil {
		return "", err
	}
	fan, err := c.StateFan()
	if err != nil {
		return "", err
	}
	return fan.Mode, nil
}

// FanPower reports whether the fan relay is on.
func (e *EAF) FanPower() (bool, error) {
	c, err := e.getConn()
	if err != nil {
		return false, err
	}
	fan, err := c.StateFan()
	if err != nil {
		return false, err
	}
	return fan.Power, nil
}

// FanSpeed returns the effective fan speed (0..100).
func (e *EAF) FanSpeed() (int, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	fan, err := c.StateFan()
	if err != nil {
		return 0, err
	}
	return fan.Speed, nil
}

// FanManualSpeed returns the manual fan speed (1..100).
func (e *EAF) FanManualSpeed() (int, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	fan, err := c.StateFan()
	if err != nil {
		return 0, err
	}
	return fan.ManualSpeed, nil
}

// FanRPM returns the fan RPM, or -1 while the tachometer is not wired.
func (e *EAF) FanRPM() (int, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	fan, err := c.StateFan()
	if err != nil {
		return 0, err
	}
	if fan.RPM == nil {
		return -1, nil
	}
	return *fan.RPM, nil
}
