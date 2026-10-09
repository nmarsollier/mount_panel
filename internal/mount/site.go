package mount

import "github.com/nmarsollier/nrm_ha_go/internal/proto"

// SetSite updates and persists the site location in a single write.
func (m *Mount) SetSite(lat, lon float64, elevation int) error {
	return m.configWrite(proto.ConfigRequest{
		Lat:       &lat,
		Lon:       &lon,
		Elevation: &elevation,
	})
}

// SetSiteLatitude updates and persists the site latitude (degrees).
func (m *Mount) SetSiteLatitude(lat float64) error {
	return m.configWrite(proto.ConfigRequest{Lat: &lat})
}

// SetSiteLongitude updates and persists the site longitude (degrees, east +).
func (m *Mount) SetSiteLongitude(lon float64) error {
	return m.configWrite(proto.ConfigRequest{Lon: &lon})
}

// SetSiteElevation updates and persists the site elevation (metres).
func (m *Mount) SetSiteElevation(elev int) error {
	return m.configWrite(proto.ConfigRequest{Elevation: &elev})
}

// configWrite sends a partial config update.
func (m *Mount) configWrite(req proto.ConfigRequest) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	_, err = c.ConfigSet(req)
	return err
}
