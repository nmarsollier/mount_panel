package amb

import (
	"math"
	"time"
)

// Rolling vibration history for the UI chart. The vibration read is
// read-and-reset: StateVibration() returns what accumulated since the previous
// read, so each sample is the RMS over its ~1 s window. Status() appends one
// sample per successful read and prunes anything older than historyWindow, so
// the buffer always holds exactly the last two minutes of activity.

// historyWindow is how long a sample stays in the rolling buffer.
const historyWindow = 5 * time.Minute

// VibrationSample is one vibration reading in the rolling window.
type VibrationSample struct {
	TsMs   int64   `json:"ts_ms"`
	VelRms float64 `json:"vel_rms"` // mm/s, ISO 10816 velocity
}

// VibrationStats summarises the rolling window. A field is nil when the window
// is empty (the frontend renders nil as "—").
type VibrationStats struct {
	Count  int      `json:"count"`
	Min    *float64 `json:"min"`
	Max    *float64 `json:"max"`
	Mean   *float64 `json:"mean"`
	Stddev *float64 `json:"stddev"`
	Last   *float64 `json:"last"`
}

// VibrationHistory is the rolling window plus its summary.
type VibrationHistory struct {
	Samples []VibrationSample `json:"samples"`
	Stats   VibrationStats    `json:"stats"`
}

// historyAppend adds a sample and drops samples older than historyWindow.
func (t *Amb) historyAppend(s VibrationSample) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := s.TsMs - historyWindow.Milliseconds()
	keep := t.history[:0]
	for _, old := range t.history {
		if old.TsMs >= cutoff {
			keep = append(keep, old)
		}
	}
	keep = append(keep, s)
	t.history = keep
}

// VibrationHistory returns a copy of the rolling window and its summary.
func (t *Amb) VibrationHistory() VibrationHistory {
	t.mu.RLock()
	samples := append([]VibrationSample(nil), t.history...)
	t.mu.RUnlock()

	return VibrationHistory{
		Samples: samples,
		Stats:   vibrationStats(samples),
	}
}

// vibrationStats computes min/max/mean/stddev/last over the samples. It returns
// zero stats with nil fields when there are no samples.
func vibrationStats(samples []VibrationSample) VibrationStats {
	if len(samples) == 0 {
		return VibrationStats{}
	}

	var sum, min, max float64
	min = samples[0].VelRms
	max = samples[0].VelRms
	for _, s := range samples {
		v := s.VelRms
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	mean := sum / float64(len(samples))

	var sq float64
	for _, s := range samples {
		d := s.VelRms - mean
		sq += d * d
	}
	stddev := math.Sqrt(sq / float64(len(samples)))
	last := samples[len(samples)-1].VelRms

	return VibrationStats{
		Count:  len(samples),
		Min:    &min,
		Max:    &max,
		Mean:   &mean,
		Stddev: &stddev,
		Last:   &last,
	}
}
