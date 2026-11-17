package mount

import "time"

// Rolling mount-movement history for the UI chart overlay. Status() is polled
// at 1 Hz by the UI; each poll appends one sample noting whether the mount is
// moving (slewing or moving an axis) and prunes anything older than
// movementWindow, mirroring the ambient vibration history in internal/amb.

// movementWindow is how long a sample stays in the rolling buffer.
const movementWindow = 5 * time.Minute

// MovementSample is one poll of the mount movement state.
type MovementSample struct {
	TsMs   int64 `json:"ts_ms"`
	Moving bool  `json:"moving"`
}

// MovementInterval is a run of consecutive moving samples, compressed for the
// chart overlay.
type MovementInterval struct {
	StartMs int64 `json:"start_ms"`
	EndMs   int64 `json:"end_ms"`
}

// MovementHistory is the compressed moving intervals over the rolling window.
type MovementHistory struct {
	Intervals []MovementInterval `json:"intervals"`
}

// movementAppend adds a sample and drops samples older than movementWindow.
func (m *Mount) movementAppend(s MovementSample) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := s.TsMs - movementWindow.Milliseconds()
	keep := m.movement[:0]
	for _, old := range m.movement {
		if old.TsMs >= cutoff {
			keep = append(keep, old)
		}
	}
	keep = append(keep, s)
	m.movement = keep
}

// MovementHistory returns the moving intervals over the rolling window.
func (m *Mount) MovementHistory() MovementHistory {
	m.mu.RLock()
	samples := append([]MovementSample(nil), m.movement...)
	m.mu.RUnlock()

	return MovementHistory{Intervals: movementIntervals(samples)}
}

// movementIntervals compresses consecutive moving samples into intervals. An
// explicit stop (a non-moving sample) breaks a run; a missed poll (no sample
// between two moving ones) does not, so a hole stays part of the same run.
func movementIntervals(samples []MovementSample) []MovementInterval {
	var intervals []MovementInterval
	prevMoving := false
	for _, s := range samples {
		if s.Moving && prevMoving {
			intervals[len(intervals)-1].EndMs = s.TsMs
		} else if s.Moving {
			intervals = append(intervals, MovementInterval{StartMs: s.TsMs, EndMs: s.TsMs})
		}
		prevMoving = s.Moving
	}
	return intervals
}
