package eaf

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

// FanSpeed returns the fan PWM duty in percent.
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
