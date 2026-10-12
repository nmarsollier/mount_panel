package mount

import "github.com/nmarsollier/nrm_ha_go/internal/mountproto"

// SetSite updates and persists the site location in a single write.
func (m *Mount) SetSite(lat, lon float64, elevation int) error {
	return m.configWrite(mountproto.ConfigRequest{
		Lat:       &lat,
		Lon:       &lon,
		Elevation: &elevation,
	})
}

// SetSiteLatitude updates and persists the site latitude (degrees).
func (m *Mount) SetSiteLatitude(lat float64) error {
	return m.configWrite(mountproto.ConfigRequest{Lat: &lat})
}

// SetSiteLongitude updates and persists the site longitude (degrees, east +).
func (m *Mount) SetSiteLongitude(lon float64) error {
	return m.configWrite(mountproto.ConfigRequest{Lon: &lon})
}

// SetSiteElevation updates and persists the site elevation (metres).
func (m *Mount) SetSiteElevation(elev int) error {
	return m.configWrite(mountproto.ConfigRequest{Elevation: &elev})
}

// configWrite sends a partial config update.
func (m *Mount) configWrite(req mountproto.ConfigRequest) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	_, err = c.ConfigSet(req)
	return err
}
