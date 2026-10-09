package proto

import (
	"encoding/binary"
	"errors"
	"time"
)

// maxFrame is the maximum JSON body length of a single frame (plus the 4-byte
// length prefix), matching the firmware's PROTO_MAX_FRAME.
const maxFrame = 2048

// frameTimeout is how long a plausible frame may stay incomplete before the
// stream is considered corrupt. A full 2048-byte body takes ~180 ms at 115200
// baud, so this is far more than a healthy link needs.
const frameTimeout = 5 * time.Second

// errFrameTimeout reports that a partial frame did not complete in time.
var errFrameTimeout = errors.New("serial frame incomplete for too long")

// encodeFrame wraps a JSON body in the wire format: [u32 little-endian
// length][JSON UTF-8].
func encodeFrame(body []byte) []byte {
	out := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(out, uint32(len(body)))
	copy(out[4:], body)
	return out
}

// frameScanner reassembles length-prefixed frames from a raw byte stream,
// dropping bytes to resynchronise when a length prefix is not plausible or the
// body is not a JSON object. Mirrors the firmware's proto_feed resync logic.
type frameScanner struct {
	buf          []byte
	pendingSince time.Time
	deadline     time.Duration // 0 = frameTimeout
}

// push feeds received bytes and invokes onFrame for each complete frame. It
// returns errFrameTimeout when a plausible frame stays incomplete past the
// deadline, so the caller can treat the stream as corrupt and reconnect.
func (s *frameScanner) push(data []byte, onFrame func([]byte)) error {
	deadline := s.deadline
	if deadline <= 0 {
		deadline = frameTimeout
	}

	s.buf = append(s.buf, data...)

	for len(s.buf) >= 4 {
		bodyLen := binary.LittleEndian.Uint32(s.buf)
		if bodyLen == 0 || bodyLen > maxFrame {
			// Not a plausible length — drop a byte and rescan.
			s.buf = s.buf[1:]
			s.pendingSince = time.Time{}
			continue
		}
		if len(s.buf) >= 5 && s.buf[4] != '{' {
			// A valid body always starts with '{'. Reject the prefix as soon
			// as its first body byte disagrees, before it can swallow later
			// frames.
			s.buf = s.buf[1:]
			s.pendingSince = time.Time{}
			continue
		}
		if len(s.buf) < 4+int(bodyLen) {
			// Frame incomplete. Track how long it has been pending.
			if s.pendingSince.IsZero() {
				s.pendingSince = time.Now()
			} else if time.Since(s.pendingSince) > deadline {
				return errFrameTimeout
			}
			return nil // wait for more bytes
		}

		body := s.buf[4 : 4+bodyLen]
		if len(body) > 0 && body[len(body)-1] == '}' {
			onFrame(body)
		}

		s.buf = s.buf[4+int(bodyLen):]
		s.pendingSince = time.Time{}
	}
	return nil
}
