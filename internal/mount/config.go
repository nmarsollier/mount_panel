package mount

// SiteLatitude returns the configured site latitude in degrees.
func (m *Mount) SiteLatitude() (float64, error) {
	c, err := m.fetchConfig()
	if err != nil {
		return 0, err
	}
	return c.Config.Lat, nil
}

// SiteLongitude returns the configured site longitude in degrees (east +).
func (m *Mount) SiteLongitude() (float64, error) {
	c, err := m.fetchConfig()
	if err != nil {
		return 0, err
	}
	return c.Config.Lon, nil
}

// SiteElevation returns the configured site elevation in metres.
func (m *Mount) SiteElevation() (int, error) {
	c, err := m.fetchConfig()
	if err != nil {
		return 0, err
	}
	return c.Config.Elevation, nil
}

// GuideRateRA returns the RA guide rate in deg/s.
func (m *Mount) GuideRateRA() (float64, error) {
	c, err := m.fetchConfig()
	if err != nil {
		return 0, err
	}
	return c.Config.GuideRateRa, nil
}

// GuideRateDec returns the DEC guide rate in deg/s.
func (m *Mount) GuideRateDec() (float64, error) {
	c, err := m.fetchConfig()
	if err != nil {
		return 0, err
	}
	return c.Config.GuideRateDec, nil
}
