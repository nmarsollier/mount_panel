package mountproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
	"testing"
	"time"

	"go.bug.st/serial"

	"github.com/nmarsollier/mount_panel/internal/proto"
)

// scriptedPort returns one scripted response frame after the first Write, then
// blocks until Close, so a single request/response round trip can be exercised
// without a real device.
type scriptedPort struct {
	response  []byte
	written   []byte
	mu        sync.Mutex
	respCh    chan struct{}
	closeCh   chan struct{}
	closeOnce sync.Once
}

func newScriptedPort(response string) *scriptedPort {
	return &scriptedPort{
		response: frame(response),
		respCh:   make(chan struct{}, 1),
		closeCh:  make(chan struct{}),
	}
}

func (p *scriptedPort) Read(buf []byte) (int, error) {
	select {
	case <-p.respCh:
		return copy(buf, p.response), nil
	case <-p.closeCh:
		return 0, errors.New("port closed")
	}
}

func (p *scriptedPort) Write(b []byte) (int, error) {
	p.mu.Lock()
	p.written = append([]byte(nil), b...)
	p.mu.Unlock()
	select {
	case p.respCh <- struct{}{}:
	default:
	}
	return len(b), nil
}

func (p *scriptedPort) Close() error {
	p.closeOnce.Do(func() { close(p.closeCh) })
	return nil
}

func (p *scriptedPort) SetMode(*serial.Mode) error                    { return nil }
func (p *scriptedPort) Drain() error                                  { return nil }
func (p *scriptedPort) ResetInputBuffer() error                       { return nil }
func (p *scriptedPort) ResetOutputBuffer() error                      { return nil }
func (p *scriptedPort) SetDTR(bool) error                             { return nil }
func (p *scriptedPort) SetRTS(bool) error                             { return nil }
func (p *scriptedPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return nil, nil
}
func (p *scriptedPort) SetReadTimeout(time.Duration) error { return nil }
func (p *scriptedPort) Break(time.Duration) error          { return nil }

func frame(body string) []byte {
	b := []byte(body)
	out := make([]byte, 4+len(b))
	binary.LittleEndian.PutUint32(out, uint32(len(b)))
	copy(out[4:], b)
	return out
}

func newQueryClient(t *testing.T, response string) (*Client, *scriptedPort) {
	t.Helper()
	p := newScriptedPort(response)
	return &Client{Client: proto.NewClient(p, "fake", time.Second)}, p
}

func TestDestinationSideOfPierEast(t *testing.T) {
	c, _ := newQueryClient(t, `{"type":"query","id":1,"ok":true,"side":"east"}`)
	defer c.Close()

	resp, err := c.DestinationSideOfPier(12.5, -35.0)
	if err != nil {
		t.Fatalf("DestinationSideOfPier: %v", err)
	}
	if resp.Side != "east" {
		t.Fatalf("side = %q, want east", resp.Side)
	}
}

func TestDestinationSideOfPierWest(t *testing.T) {
	c, _ := newQueryClient(t, `{"type":"query","id":1,"ok":true,"side":"west"}`)
	defer c.Close()

	resp, err := c.DestinationSideOfPier(12.5, -35.0)
	if err != nil {
		t.Fatalf("DestinationSideOfPier: %v", err)
	}
	if resp.Side != "west" {
		t.Fatalf("side = %q, want west", resp.Side)
	}
}

func TestDestinationSideOfPierUnreachable(t *testing.T) {
	c, _ := newQueryClient(t, `{"type":"query","id":1,"ok":false,"error":"Target unreachable"}`)
	defer c.Close()

	if _, err := c.DestinationSideOfPier(12.5, 60.0); err == nil {
		t.Fatal("expected an error for an unreachable destination")
	}
}

func TestDestinationSideOfPierMalformedResponse(t *testing.T) {
	c, _ := newQueryClient(t, `not-json`)
	defer c.Close()

	if _, err := c.DestinationSideOfPier(12.5, -35.0); err == nil {
		t.Fatal("expected an error for a malformed response")
	}
}

func TestDestinationSideOfPierRequestFrame(t *testing.T) {
	c, p := newQueryClient(t, `{"type":"query","id":1,"ok":true,"side":"east"}`)
	defer c.Close()

	if _, err := c.DestinationSideOfPier(12.5, -35.0); err != nil {
		t.Fatalf("DestinationSideOfPier: %v", err)
	}

	p.mu.Lock()
	written := append([]byte(nil), p.written...)
	p.mu.Unlock()

	if len(written) < 4 {
		t.Fatalf("request frame too short: %d bytes", len(written))
	}
	body := written[4:]
	for _, want := range []string{`"type":"query"`, `"query":"destination_side_of_pier"`, `"ra":12.5`, `"dec":-35`} {
		if !bytes.Contains(body, []byte(want)) {
			t.Errorf("request frame missing %s: %s", want, body)
		}
	}
}
