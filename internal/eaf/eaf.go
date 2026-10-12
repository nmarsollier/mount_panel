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

	"github.com/nmarsollier/nrm_ha_go/internal/eafproto"
	"github.com/nmarsollier/nrm_ha_go/internal/tools/logx"
)

// errNotConnected is returned when an operation is requested but the serial
// link to the accessory controller is down.
var errNotConnected = errors.New("accessory controller not connected")

// retryInterval is how long to wait before re-attempting a failed connection.
const retryInterval = 2 * time.Second

// EAF is the model of the accessory controller: it manages the serial
// connection and exposes the focuser and fan operations to the UI.
type EAF struct {
	mu sync.RWMutex

	conn       *eafproto.Client // nil = disconnected
	serialPort string

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

// Init starts the accessory model and its background connection loop. It does
// not block: the controller is connected whenever it becomes available.
func Init(serialPort string) {
	once.Do(func() {
		inst = &EAF{serialPort: serialPort}
		go inst.connectLoop()
	})
}

// connectLoop keeps the serial link up: it connects, waits for the link to
// drop, then reconnects.
func (e *EAF) connectLoop() {
	logx.Info(eafCtx(), "connect loop started")
	for {
		c, err := e.tryConnect()
		if err != nil {
			logx.Warn(eafCtx(), "connect failed: "+err.Error())
			time.Sleep(retryInterval)
			continue
		}

		<-c.Done() // block until the link drops
		reason := ""
		if err := c.Err(); err != nil {
			reason = ": " + err.Error()
		}
		logx.Warn(eafCtx(), "accessory disconnected"+reason)
		e.disconnect(c)
	}
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
