package proto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial"
)

// DefaultTimeout is the round-trip timeout once a device is connected.
const DefaultTimeout = 5 * time.Second

var errTimeout = errors.New("serial request timed out")

// Client owns a serial port and correlates length-prefixed JSON responses to
// the request that asked for them. It is transport-only: it knows nothing about
// the mount or accessory message contracts, only that every response carries an
// `id` at its top level. A single request is in flight at a time.
type Client struct {
	port serial.Port
	name string

	reqMu     sync.Mutex // serialises in-flight requests (one at a time)
	pendingMu sync.Mutex // guards the pending map (reader + requester)
	pending   map[uint32]chan json.RawMessage
	nextID    uint32
	timeout   time.Duration

	failMu  sync.Mutex // guards failErr (first reason the link dropped)
	failErr error

	closeOnce sync.Once   // closes the port exactly once
	closed    atomic.Bool // set before the port is closed
	closeErr  error       // port.Close result
	done      chan struct{}
}

// Open connects to the given serial port and starts the reader goroutine. It
// claims the port exclusively; a second Open (or detector) on the same port
// fails until the first client closes.
func Open(name string) (*Client, error) {
	if !ClaimPort(name) {
		return nil, errors.New("serial port already in use")
	}
	p, err := OpenPort(name)
	if err != nil {
		ReleasePort(name)
		return nil, err
	}
	return NewClient(p, name, DefaultTimeout), nil
}

// NewClient builds a Client on an open port and starts its reader goroutine.
func NewClient(p serial.Port, name string, timeout time.Duration) *Client {
	c := &Client{
		port:    p,
		name:    name,
		pending: map[uint32]chan json.RawMessage{},
		timeout: timeout,
		done:    make(chan struct{}),
	}
	go c.readLoop()
	return c
}

// PortName returns the serial port name this client is connected to.
func (c *Client) PortName() string {
	return c.name
}

// Done is closed when the connection is lost or the client is closed.
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// SetTimeout changes the round-trip timeout, used after a detect handshake.
func (c *Client) SetTimeout(d time.Duration) {
	c.timeout = d
}

// NextID allocates a fresh request id for this connection. Ids start at 1; 0 is
// reserved for events.
func (c *Client) NextID() uint32 {
	return atomic.AddUint32(&c.nextID, 1)
}

// Close releases the serial port. It is idempotent and does not acquire the
// request mutex, so a failing RoundTrip can call it to interrupt a blocked
// read or write without deadlocking on reqMu. It also closes Done so a blocked
// connection loop wakes even when the reader goroutine is stuck in Read.
func (c *Client) Close() error {
	c.closed.Store(true)
	c.closeOnce.Do(func() {
		c.closeErr = c.port.Close()
		ReleasePort(c.name)
		close(c.done)
	})
	return c.closeErr
}

// readLoop reads frames from the serial port and dispatches them to the
// pending request waiting on that id. Every exit path goes through fail, which
// closes the client and Done, so the owner can reconnect.
func (c *Client) readLoop() {
	buf := make([]byte, 4096)
	var scanner frameScanner

	for {
		n, err := c.port.Read(buf)
		if n > 0 {
			if serr := scanner.push(buf[:n], c.dispatch); serr != nil {
				c.fail(serr)
				return
			}
		}
		if err != nil {
			c.fail(err)
			return
		}
	}
}

// dispatch routes a response body to the request waiting on its id.
func (c *Client) dispatch(body []byte) {
	var id struct {
		ID uint32 `json:"id"`
	}
	if err := json.Unmarshal(body, &id); err != nil {
		return
	}

	c.pendingMu.Lock()
	ch, ok := c.pending[id.ID]
	c.pendingMu.Unlock()

	if ok {
		ch <- body
	}
}

// failAll wakes every pending request with an error so callers do not hang on
// a lost serial link. Sends are non-blocking to avoid stalling the reader.
func (c *Client) failAll(err error) {
	c.pendingMu.Lock()
	chans := make([]chan json.RawMessage, 0, len(c.pending))
	for _, ch := range c.pending {
		chans = append(chans, ch)
	}
	c.pendingMu.Unlock()

	msg := json.RawMessage(fmt.Sprintf(`{"ok":false,"error":%q}`, err.Error()))
	for _, ch := range chans {
		select {
		case ch <- msg:
		default:
		}
	}
}

// fail invalidates the connection: it wakes every pending request and closes
// the port, which unblocks the reader and tears the session down so the owner
// can reconnect. It does not acquire reqMu, so RoundTrip may call it safely.
func (c *Client) fail(err error) {
	c.failMu.Lock()
	if c.failErr == nil {
		c.failErr = err
	}
	c.failMu.Unlock()

	c.failAll(err)
	c.Close()
}

// Err returns the first error that tore the connection down, or nil when the
// client was closed cleanly.
func (c *Client) Err() error {
	c.failMu.Lock()
	defer c.failMu.Unlock()
	return c.failErr
}

// RoundTrip sends a framed request body and waits for the response with the
// matching id.
func (c *Client) RoundTrip(ctx context.Context, id uint32, body []byte) (json.RawMessage, error) {
	c.reqMu.Lock()
	defer c.reqMu.Unlock()

	if c.closed.Load() {
		return nil, errors.New("serial connection closed")
	}
	if len(body) > maxFrame {
		return nil, errors.New("request frame too large")
	}

	ch := make(chan json.RawMessage, 1)
	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	if _, err := c.port.Write(encodeFrame(body)); err != nil {
		c.fail(err)
		return nil, err
	}

	select {
	case raw := <-ch:
		return raw, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(c.timeout):
		c.fail(errTimeout)
		return nil, errTimeout
	}
}
