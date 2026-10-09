package proto

import (
	"errors"
	"sync"
	"testing"
	"time"

	"go.bug.st/serial"
)

// fakePort is a scripted serial.Port for testing the client without a device.
// Read blocks until Close; Write can be made to block until Close (writeBlock)
// to exercise interrupting a blocked write.
type fakePort struct {
	closeCh    chan struct{}
	closeOnce  sync.Once
	writeBlock bool
	writeStart chan struct{}
	writeErr   error
}

func newFakePort() *fakePort {
	return &fakePort{closeCh: make(chan struct{}), writeStart: make(chan struct{})}
}

func (p *fakePort) Read(_ []byte) (int, error) {
	<-p.closeCh
	return 0, errors.New("port closed")
}

func (p *fakePort) Write(b []byte) (int, error) {
	if p.writeBlock {
		close(p.writeStart)
		<-p.closeCh
		return 0, errors.New("port closed")
	}
	if p.writeErr != nil {
		return 0, p.writeErr
	}
	return len(b), nil
}

func (p *fakePort) Close() error {
	p.closeOnce.Do(func() { close(p.closeCh) })
	return nil
}

func (p *fakePort) SetMode(*serial.Mode) error { return nil }
func (p *fakePort) Drain() error               { return nil }
func (p *fakePort) ResetInputBuffer() error    { return nil }
func (p *fakePort) ResetOutputBuffer() error   { return nil }
func (p *fakePort) SetDTR(bool) error          { return nil }
func (p *fakePort) SetRTS(bool) error          { return nil }
func (p *fakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return nil, nil
}
func (p *fakePort) SetReadTimeout(time.Duration) error { return nil }
func (p *fakePort) Break(time.Duration) error          { return nil }

func TestRoundTripTimeoutFailsClient(t *testing.T) {
	c := newClient(newFakePort(), "fake", 20*time.Millisecond)

	_, err := c.State()
	if !errors.Is(err, errTimeout) {
		t.Fatalf("want errTimeout, got %v", err)
	}
	if !c.closed.Load() {
		t.Fatal("client not marked closed after timeout")
	}
	select {
	case <-c.Done():
	case <-time.After(time.Second):
		t.Fatal("Done not closed after timeout")
	}
}

func TestCloseInterruptsBlockedWrite(t *testing.T) {
	p := newFakePort()
	p.writeBlock = true
	c := newClient(p, "fake", time.Second)

	errCh := make(chan error, 1)
	go func() {
		_, err := c.State()
		errCh <- err
	}()

	<-p.writeStart // wait until the write is blocked

	c.Close()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected an error from the interrupted write")
		}
	case <-time.After(time.Second):
		t.Fatal("roundTrip did not return after Close")
	}
}
