package eaf

// SetFanMode selects the fan mode: "normal", "soft" or "off".
func (e *EAF) SetFanMode(mode string) error {
	c, err := e.getConn()
	if err != nil {
		return err
	}
	_, err = c.FanSetMode(mode)
	return err
}
