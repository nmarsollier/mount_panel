// Package eaf models the accessory controller (focuser, fan, environment) that
// hangs off the NRM-HA mount over a second USB CDC-ACM link. It mirrors the
// mount package: a singleton that keeps the serial link up and translates
// high-level operations into eafproto calls. No state cache: every read
// fetches a fresh snapshot over the serial link.
package eaf

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nmarsollier/mount_panel/internal/eafproto"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// errNotConnected is returned when an operation is requested but the serial
// link to the accessory controller is down.
var errNotConnected = errors.New("accessory controller not connected")

// retryInterval is the delay between startup connection attempts while the USB
// device enumerates.
const retryInterval = 2 * time.Second

// startupAttempts bounds the startup retry loop so it does not spin forever
// when the accessory never appears; reconnection is then manual, from the UI.
const startupAttempts = 15

// EAF is the model of the accessory controller: it manages the serial
// connection and exposes the focuser and fan operations to the UI.
type EAF struct {
	mu sync.RWMutex

	conn       *eafproto.Client // nil = disconnected
	serialPort string
	connecting bool

	// Virtual absolute position (Alpaca, int32 non-negative) and whether a
	// reference has been established. Reset when the controller reconnects.
	position      int32
	positionValid bool

	// averagePeriod is the ObservingConditions averaging window, in hours
	// (0 = instantaneous). Stored in Go; it does not change the sensor period.
	averagePeriod float64

	connected bool // logical client link (Alpaca Connected)
}

var (
	once sync.Once
	inst *EAF
)

// Get returns the singleton accessory model. It is non-nil after Init.
func Get() *EAF {
	return inst
}

// Init starts the accessory model. It makes a single startup connection
// attempt; further attempts are manual, from the UI, via Connect.
func Init(serialPort string) {
	once.Do(func() {
		inst = &EAF{serialPort: serialPort}
		go inst.connectAtStartup()
	})
}

// connectAtStartup retries the initial connection until the accessory answers
// or the startup window elapses. The USB CDC port enumerates asynchronously
// after launch, so the first attempt can race it; this covers that gap. After
// the link is up, the connection is managed manually (Connect/watch), so an
// unplug is not reconnected in the background.
func (e *EAF) connectAtStartup() {
	for i := 0; i < startupAttempts; i++ {
		if e.Connect() == nil && e.IsConnected() {
			return
		}
		time.Sleep(retryInterval)
	}
}

// Connect makes a single attempt to open and verify the accessory controller.
// It is a no-op when the link is already up or a connection attempt is in
// flight, and returns the error when the attempt fails.
func (e *EAF) Connect() error {
	e.mu.Lock()
	if e.conn != nil || e.connecting {
		e.mu.Unlock()
		return nil
	}
	e.connecting = true
	e.mu.Unlock()
	defer func() {
		e.mu.Lock()
		e.connecting = false
		e.mu.Unlock()
	}()

	c, err := e.tryConnect()
	if err != nil {
		logx.Warn(eafCtx(), "connect failed: "+err.Error())
		return err
	}
	go e.watch(c)
	return nil
}

// watch tears the link down when the reader notices it has dropped. The link is
// not reconnected automatically: the user re-connects from the UI.
func (e *EAF) watch(c *eafproto.Client) {
	<-c.Done()
	reason := ""
	if err := c.Err(); err != nil {
		reason = ": " + err.Error()
	}
	logx.Warn(eafCtx(), "accessory disconnected"+reason)
	e.disconnect(c)
}

// tryConnect opens and verifies the accessory controller.
func (e *EAF) tryConnect() (*eafproto.Client, error) {
	var (
		c   *eafproto.Client
		err error
	)
	if e.serialPort != "" {
		c, err = eafproto.Open(e.serialPort)
	} else {
		c, err = eafproto.Detect()
	}
	if err != nil {
		return nil, err
	}

	e.mu.Lock()
	e.conn = c
	e.position = 0
	e.positionValid = true
	e.mu.Unlock()

	logx.Info(eafCtx(), "connected on "+c.PortName()+" boot="+c.BootID())
	return c, nil
}

// disconnect closes the client and clears the connection.
func (e *EAF) disconnect(c *eafproto.Client) {
	c.Close()
	e.mu.Lock()
	if e.conn == c {
		e.conn = nil
	}
	e.positionValid = false
	e.mu.Unlock()
}

// getConn returns the current serial connection, or an error when disconnected.
func (e *EAF) getConn() (*eafproto.Client, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.conn == nil {
		return nil, errNotConnected
	}
	return e.conn, nil
}

// IsConnected reports whether the serial link to the accessory controller is up.
func (e *EAF) IsConnected() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.conn != nil
}

func eafCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "EAF"))
}
