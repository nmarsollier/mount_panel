package mount

// DeviceState is a coherent snapshot of the Alpaca DeviceState properties,
// built from a single state read (plus one config read for altitude/azimuth).
type DeviceState struct {
	Altitude       float64
	AtHome         bool
	AtPark         bool
	Azimuth        float64
	Declination    float64
	IsPulseGuiding bool
	RightAscension float64
	SideOfPier     int
	Slewing        bool
	Tracking       bool
}

// DeviceState reads the mount once and returns the operational properties the
// Alpaca DeviceState endpoint aggregates, avoiding one serial round-trip per
// property and a mix of different instants.
func (m *Mount) DeviceState() (DeviceState, error) {
	s, err := m.fetchState()
	if err != nil {
		return DeviceState{}, err
	}
	lat, err := m.SiteLatitude()
	if err != nil {
		return DeviceState{}, err
	}
	alt, az := horizontal(s.Ra, s.Dec, lat, s.Lst)

	return DeviceState{
		Altitude:       alt,
		AtHome:         s.AtHome,
		AtPark:         s.AtPark,
		Azimuth:        az,
		Declination:    s.Dec,
		IsPulseGuiding: s.Guiding,
		RightAscension: s.Ra,
		SideOfPier:     pierSideState(s),
		Slewing:        slewingState(s),
		Tracking:       trackingState(s),
	}, nil
}
