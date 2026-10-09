package proto

import "encoding/binary"

// maxFrame is the maximum JSON body length of a single frame (plus the 4-byte
// length prefix), matching the firmware's PROTO_MAX_FRAME.
const maxFrame = 2048

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
	buf []byte
}

// push feeds received bytes and invokes onFrame for each complete frame.
func (s *frameScanner) push(data []byte, onFrame func([]byte)) {
	s.buf = append(s.buf, data...)

	for len(s.buf) >= 4 {
		bodyLen := binary.LittleEndian.Uint32(s.buf)
		if bodyLen == 0 || bodyLen > maxFrame {
			// Not a plausible length — drop a byte and rescan.
			s.buf = s.buf[1:]
			continue
		}
		if len(s.buf) < 4+int(bodyLen) {
			return // frame incomplete, wait for more bytes
		}

		body := s.buf[4 : 4+bodyLen]
		if len(body) > 0 && body[0] == '{' && body[len(body)-1] == '}' {
			onFrame(body)
		}

		s.buf = s.buf[4+int(bodyLen):]
	}
}
