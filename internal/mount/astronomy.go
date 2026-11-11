package mount

import (
	"math"
	"time"
)

// unixToJD converts a UTC time to a Julian Date.
func unixToJD(t time.Time) float64 {
	return float64(t.Unix())/86400.0 + 2440587.5
}

// gmstHours computes Greenwich Mean Sidereal Time in hours (USNO approximation).
func gmstHours(jd float64) float64 {
	T := (jd - 2451545.0) / 36525.0
	gmst := 280.46061837 + 360.98564736629*(jd-2451545.0) +
		0.000387933*T*T - (T*T*T)/38710000.0
	gmst = math.Mod(gmst, 360.0)
	if gmst < 0 {
		gmst += 360.0
	}
	return gmst / 15.0
}

// lstHours computes Local Mean Sidereal Time in hours from the host clock and a
// site longitude. Used only as a fallback when the mount is disconnected.
func lstHours(lonDeg float64, t time.Time) float64 {
	lst := gmstHours(unixToJD(t)) + lonDeg/15.0
	lst = math.Mod(lst, 24.0)
	if lst < 0 {
		lst += 24.0
	}
	return lst
}

// horizontal computes altitude and azimuth (degrees) from an equatorial
// position, site latitude and the local sidereal time.
func horizontal(raHours, decDeg, latDeg, lst float64) (alt, az float64) {
	haRad := (lst - raHours) * 15.0 * math.Pi / 180.0
	decRad := decDeg * math.Pi / 180.0
	latRad := latDeg * math.Pi / 180.0

	altRad := math.Asin(math.Sin(decRad)*math.Sin(latRad) +
		math.Cos(decRad)*math.Cos(latRad)*math.Cos(haRad))
	alt = altRad * 180.0 / math.Pi

	cosAz := (math.Sin(decRad) - math.Sin(altRad)*math.Sin(latRad)) /
		(math.Cos(altRad) * math.Cos(latRad))
	cosAz = math.Max(-1, math.Min(1, cosAz))
	azRad := math.Acos(cosAz)
	if math.Sin(haRad) > 0 {
		azRad = 2*math.Pi - azRad
	}
	az = azRad * 180.0 / math.Pi
	return alt, az
}

// Altitude returns the current telescope altitude in degrees.
func (m *Mount) Altitude() (float64, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	lat, err := m.SiteLatitude()
	if err != nil {
		return 0, err
	}
	alt, _ := horizontal(s.Ra, s.Dec, lat, s.Lst)
	return alt, nil
}

// Azimuth returns the current telescope azimuth in degrees (0 = North).
func (m *Mount) Azimuth() (float64, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	lat, err := m.SiteLatitude()
	if err != nil {
		return 0, err
	}
	_, az := horizontal(s.Ra, s.Dec, lat, s.Lst)
	return az, nil
}

// SiderealTime returns the local mean sidereal time in hours as reported by the
// mount. It does not fall back to the host clock: a communication failure is
// returned as an error rather than masked as a valid reading.
func (m *Mount) SiderealTime() (float64, error) {
	s, err := m.fetchState()
	if err != nil {
		return 0, err
	}
	return s.Lst, nil
}
