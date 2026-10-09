package mount

import (
	"fmt"
	"math"
	"time"
)

// Status is a coherent snapshot of the mount formatted for the UI, mirroring
// the firmware's GET /api/status shape.
type Status struct {
	Status   string `json:"status"`
	Tracking string `json:"tracking"`
	Power    bool   `json:"power"`
	Ra       string `json:"ra"`
	Dec      string `json:"dec"`
	Lst      string `json:"lst"`
	PierSide string `json:"pier_side"`
	Time     string `json:"time"`
	IsHome   bool   `json:"is_home"`
	Settings struct {
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Elevation int     `json:"elevation"`
	} `json:"settings"`
	Debug struct {
		RaAxisDeg  float64 `json:"ra_axis_deg"`
		DecAxisDeg float64 `json:"dec_axis_deg"`
		RaSteps    int64   `json:"ra_steps"`
		DecSteps   int64   `json:"dec_steps"`
		RaSpeed    float64 `json:"ra_speed"`
		DecSpeed   float64 `json:"dec_speed"`
		Guiding    bool    `json:"guiding"`
		Microsteps int     `json:"microsteps"`
		Limits     struct {
			RaMin  float64 `json:"ra_min"`
			RaMax  float64 `json:"ra_max"`
			DecMin float64 `json:"dec_min"`
			DecMax float64 `json:"dec_max"`
		} `json:"limits"`
		UptimeS uint64 `json:"uptime_s"`
	} `json:"debug"`
}

// Status builds the UI snapshot from a fresh state/config read and the host
// clock.
func (m *Mount) Status() Status {
	s := m.fetchState()
	cfg := m.fetchConfig()

	st := Status{
		Status:   "disconnected",
		Tracking: "none",
		Power:    false,
		Ra:       "--:--:--",
		Dec:      "--:--:--",
		Lst:      formatHours(m.SiderealTime()),
		PierSide: "East",
		Time:     time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
	st.Debug.Microsteps = 64 // MOTORS_MICROSTEPS (firmware build constant)

	if s != nil {
		st.Status = s.State
		st.Tracking = s.Tracking
		st.Power = s.Power
		st.Ra = formatHours(s.Ra)
		st.Dec = formatDec(s.Dec)
		st.IsHome = s.AtHome
		if s.PierSide == "west" {
			st.PierSide = "West"
		} else {
			st.PierSide = "East"
		}
		st.Debug.RaAxisDeg = s.RaAxisDeg
		st.Debug.DecAxisDeg = s.DecAxisDeg
		st.Debug.RaSteps = s.RaSteps
		st.Debug.DecSteps = s.DecSteps
		st.Debug.RaSpeed = s.RaSpeedDps
		st.Debug.DecSpeed = s.DecSpeedDps
		st.Debug.Guiding = s.Guiding
		st.Debug.Limits.RaMin = s.Limits.RaMin
		st.Debug.Limits.RaMax = s.Limits.RaMax
		st.Debug.Limits.DecMin = s.Limits.DecMin
		st.Debug.Limits.DecMax = s.Limits.DecMax
		st.Debug.UptimeS = s.UptimeMs / 1000
	}

	if cfg != nil {
		st.Settings.Lat = cfg.Config.Lat
		st.Settings.Lon = cfg.Config.Lon
		st.Settings.Elevation = cfg.Config.Elevation
	}

	return st
}

// formatHours renders decimal hours as "HH:MM:SS.ss".
func formatHours(h float64) string {
	total := math.Abs(h)
	hh := int(total)
	rem := (total - float64(hh)) * 60
	mm := int(rem)
	ss := (rem - float64(mm)) * 60
	return fmt.Sprintf("%02d:%02d:%05.2f", hh, mm, ss)
}

// formatDec renders decimal degrees as "±DD:MM:SS.ss".
func formatDec(d float64) string {
	sign := "+"
	if d < 0 {
		sign = "-"
	}
	total := math.Abs(d)
	dd := int(total)
	rem := (total - float64(dd)) * 60
	mm := int(rem)
	ss := (rem - float64(mm)) * 60
	return fmt.Sprintf("%s%02d:%02d:%05.2f", sign, dd, mm, ss)
}
