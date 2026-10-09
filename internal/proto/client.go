package proto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
)

// defaultTimeout is the round-trip timeout once a mount is connected.
const defaultTimeout = 5 * time.Second

// Client talks to the mount over the USB CDC-ACM protocol. It owns the serial
// port, a background reader that routes responses by request id, and the
// request/response bookkeeping. A single request is in flight at a time.
type Client struct {
	port serial.Port
	name string

	mu        sync.Mutex // serialises writes and in-flight requests
	pendingMu sync.Mutex // guards the pending map (reader + requester)
	pending   map[uint32]chan json.RawMessage
	nextID    uint32
	timeout   time.Duration

	done   chan struct{} // closed when the reader exits
	closed bool
}

// Open connects to the given serial port and starts the reader goroutine.
func Open(name string) (*Client, error) {
	p, err := openPort(name)
	if err != nil {
		return nil, err
	}
	return newClient(p, name, defaultTimeout), nil
}

// newClient builds a Client on an open port and starts its reader goroutine.
func newClient(p serial.Port, name string, timeout time.Duration) *Client {
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

// Close releases the serial port.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	return c.port.Close()
}

// readLoop reads frames from the serial port and dispatches them to the
// pending request waiting on that id.
func (c *Client) readLoop() {
	defer close(c.done)

	buf := make([]byte, 4096)
	var scanner frameScanner

	for {
		n, err := c.port.Read(buf)
		if n > 0 {
			scanner.push(buf[:n], c.dispatch)
		}
		if err != nil {
			c.failAll(err)
			return
		}
	}
}

// dispatch routes a response body to the request waiting on its id.
func (c *Client) dispatch(body []byte) {
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return
	}

	c.pendingMu.Lock()
	ch, ok := c.pending[env.ID]
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

// roundTrip sends a request and waits for the response with the matching id.
func (c *Client) roundTrip(ctx context.Context, reqType string, payload any) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil, errors.New("serial connection closed")
	}

	c.nextID++
	id := c.nextID

	body, err := buildRequest(reqType, id, payload)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	select {
	case raw := <-ch:
		return raw, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(c.timeout):
		return nil, errTimeout
	}
}

// request performs a round trip and unmarshals the response into out, checking
// the shared ok/error envelope first.
func (c *Client) request(ctx context.Context, reqType string, payload any, out any) error {
	raw, err := c.roundTrip(ctx, reqType, payload)
	if err != nil {
		return err
	}

	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	if !env.OK {
		return &Error{Message: env.Error}
	}

	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// buildRequest flattens a family payload into the common request envelope.
func buildRequest(reqType string, id uint32, payload any) ([]byte, error) {
	m := map[string]any{"type": reqType, "id": id}

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		var fields map[string]any
		if err := json.Unmarshal(b, &fields); err != nil {
			return nil, err
		}
		for k, v := range fields {
			m[k] = v
		}
	}

	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if len(body) > maxFrame {
		return nil, errors.New("request frame too large")
	}
	return body, nil
}
