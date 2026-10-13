package eaf

// SetFanPower switches the fan's 12 V relay (the power channel).
func (e *EAF) SetFanPower(on bool) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	return c.FanSetPower(on)
}

// SetFanSpeed sets the fan PWM duty (0..100 %).
func (e *EAF) SetFanSpeed(pct int) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	return c.FanSetSpeed(pct)
}
