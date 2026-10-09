package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
	"time"
)

func TestScannerResyncsAfterBogusLength(t *testing.T) {
	var s frameScanner

	const body = `{"type":"state","id":1,"ok":true}`
	var data []byte
	bogus := make([]byte, 4)
	binary.LittleEndian.PutUint32(bogus, 2048)
	data = append(data, bogus...)
	for i := 0; i < 10; i++ {
		data = append(data, encodeFrame([]byte(body))...)
	}

	var got [][]byte
	if err := s.push(data, func(b []byte) { got = append(got, b) }); err != nil {
		t.Fatalf("push: %v", err)
	}

	if len(got) != 10 {
		t.Fatalf("delivered %d frames, want 10", len(got))
	}
	for i, g := range got {
		if !bytes.Equal(g, []byte(body)) {
			t.Fatalf("frame %d = %s, want %s", i, g, body)
		}
	}
}

func TestScannerFrameTimeout(t *testing.T) {
	var s frameScanner
	s.deadline = 20 * time.Millisecond

	// A plausible length whose body starts with '{' but never completes.
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, 100)
	partial := append(prefix, '{')

	if err := s.push(partial, nil); err != nil {
		t.Fatalf("first push: %v", err)
	}
	time.Sleep(40 * time.Millisecond)

	if err := s.push([]byte("abc"), nil); !errors.Is(err, errFrameTimeout) {
		t.Fatalf("want errFrameTimeout, got %v", err)
	}
}
