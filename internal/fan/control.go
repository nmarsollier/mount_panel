package fan

// SetFanMode sets the fan mode: "off", "auto" or "manual".
func (f *Fan) SetFanMode(mode string) error {
	c, err := f.getConn()
	if err != nil {
		return err
	}
	return c.FanSetMode(mode)
}

// SetFanSpeed sets the manual fan speed (1..100).
func (f *Fan) SetFanSpeed(speed int) error {
	c, err := f.getConn()
	if err != nil {
		return err
	}
	return c.FanSetSpeed(speed)
}

// FanMode returns the current fan mode ("off", "auto" or "manual").
func (f *Fan) FanMode() (string, error) {
	c, err := f.getConn()
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
func (f *Fan) FanPower() (bool, error) {
	c, err := f.getConn()
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
func (f *Fan) FanSpeed() (int, error) {
	c, err := f.getConn()
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
func (f *Fan) FanManualSpeed() (int, error) {
	c, err := f.getConn()
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
func (f *Fan) FanRPM() (int, error) {
	c, err := f.getConn()
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
