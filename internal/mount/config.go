package mount

// SiteLatitude returns the configured site latitude in degrees.
func (m *Mount) SiteLatitude() float64 {
	c := m.fetchConfig()
	if c == nil {
		return 0
	}
	return c.Config.Lat
}

// SiteLongitude returns the configured site longitude in degrees (east +).
func (m *Mount) SiteLongitude() float64 {
	c := m.fetchConfig()
	if c == nil {
		return 0
	}
	return c.Config.Lon
}

// SiteElevation returns the configured site elevation in metres.
func (m *Mount) SiteElevation() int {
	c := m.fetchConfig()
	if c == nil {
		return 0
	}
	return c.Config.Elevation
}

// GuideRateRA returns the RA guide rate in deg/s.
func (m *Mount) GuideRateRA() float64 {
	c := m.fetchConfig()
	if c == nil {
		return 0
	}
	return c.Config.GuideRateRa
}

// GuideRateDec returns the DEC guide rate in deg/s.
func (m *Mount) GuideRateDec() float64 {
	c := m.fetchConfig()
	if c == nil {
		return 0
	}
	return c.Config.GuideRateDec
}
