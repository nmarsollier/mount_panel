package eaf

// FanMode returns the current fan mode ("normal", "soft" or "off").
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

// FanPWM returns the applied fan PWM duty in percent.
func (e *EAF) FanPWM() (int, error) {
	c, err := e.getConn()
	if err != nil {
		return 0, err
	}
	fan, err := c.StateFan()
	if err != nil {
		return 0, err
	}
	return fan.OutputPWMPct, nil
}

// FanModeNumber maps a fan mode to its Switch channel value (0=off, 1=soft,
// 2=normal).
func FanModeNumber(mode string) int {
	switch mode {
	case "soft":
		return 1
	case "normal":
		return 2
	default:
		return 0 // "off"
	}
}

// FanModeFromNumber maps a Switch channel value back to a fan mode.
func FanModeFromNumber(n int) string {
	switch n {
	case 1:
		return "soft"
	case 2:
		return "normal"
	default:
		return "off"
	}
}
