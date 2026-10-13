package eaf

// SetFanMode sets the fan mode: "off", "auto" or "manual".
func (e *EAF) SetFanMode(mode string) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	return c.FanSetMode(mode)
}

// SetFanSpeed sets the manual fan speed (1..100).
func (e *EAF) SetFanSpeed(speed int) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	return c.FanSetSpeed(speed)
}
