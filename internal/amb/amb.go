// Package amb models the ambient temperature/humidity and vibration module
// that hangs off the BFT mount over its own USB CDC-ACM link. It mirrors the
// eaf package: a singleton that keeps the serial link up and translates
// high-level operations into eafproto calls. No state cache: every read fetches
// a fresh snapshot.
package amb

import (
	"context"
	"errors"
	"sync"

	"github.com/nmarsollier/mount_panel/internal/eafproto"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// errNotConnected is returned when an operation is requested but the serial
// link to the ambient module is down.
var errNotConnected = errors.New("ambient module not connected")

// Amb is the model of the ambient/vibration module: it manages the serial
// connection and exposes the ambient and vibration readings to the UI.
type Amb struct {
	mu sync.RWMutex

	conn       *eafproto.Client // nil = disconnected
	serialPort string
	connecting bool

	// averagePeriod is the ObservingConditions averaging window, in hours
	// (0 = instantaneous). Stored in Go; it does not change the sensor period.
	averagePeriod float64

	connected bool // logical client link (Alpaca Connected)
}

var (
	once sync.Once
	inst *Amb
)

// Get returns the singleton ambient model. It is non-nil after Init.
func Get() *Amb {
	return inst
}

// Init starts the ambient model. It makes a single startup connection
// attempt; further attempts are manual, from the UI, via Connect.
func Init(serialPort string) {
	once.Do(func() {
		inst = &Amb{serialPort: serialPort}
		go inst.Connect()
	})
}

// Connect makes a single attempt to open and verify the ambient module. It
// is a no-op when the link is already up or a connection attempt is in flight,
// and returns the error when the attempt fails.
func (t *Amb) Connect() error {
	t.mu.Lock()
	if t.conn != nil || t.connecting {
		t.mu.Unlock()
		return nil
	}
	t.connecting = true
	t.mu.Unlock()
	defer func() {
		t.mu.Lock()
		t.connecting = false
		t.mu.Unlock()
	}()

	c, err := t.tryConnect()
	if err != nil {
		logx.Warn(ambCtx(), "connect failed: "+err.Error())
		return err
	}
	go t.watch(c)
	return nil
}

// Disconnect closes the serial link to the ambient module. It is a no-op
// when the link is already down.
func (t *Amb) Disconnect() {
	t.mu.RLock()
	c := t.conn
	t.mu.RUnlock()
	if c == nil {
		return
	}
	t.disconnect(c)
}

// watch tears the link down when the reader notices it has dropped. The link is
// not reconnected automatically: the user re-connects from the UI.
func (t *Amb) watch(c *eafproto.Client) {
	<-c.Done()
	reason := ""
	if err := c.Err(); err != nil {
		reason = ": " + err.Error()
	}
	logx.Warn(ambCtx(), "ambient module disconnected"+reason)
	t.disconnect(c)
}

// tryConnect opens and verifies the ambient module.
func (t *Amb) tryConnect() (*eafproto.Client, error) {
	var (
		c   *eafproto.Client
		err error
	)
	if t.serialPort != "" {
		c, err = eafproto.Open(t.serialPort)
	} else {
		c, err = eafproto.Detect("amb")
	}
	if err != nil {
		return nil, err
	}

	t.mu.Lock()
	t.conn = c
	t.mu.Unlock()

	logx.Info(ambCtx(), "connected on "+c.PortName()+" boot="+c.BootID())
	return c, nil
}

// disconnect closes the client and clears the connection.
func (t *Amb) disconnect(c *eafproto.Client) {
	c.Close()
	t.mu.Lock()
	if t.conn == c {
		t.conn = nil
	}
	t.mu.Unlock()
}

// getConn returns the current serial connection, or an error when disconnected.
func (t *Amb) getConn() (*eafproto.Client, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.conn == nil {
		return nil, errNotConnected
	}
	return t.conn, nil
}

// IsConnected reports whether the serial link to the ambient module is up.
func (t *Amb) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.conn != nil
}

// Connected reports the logical client link (Alpaca Connected).
func (t *Amb) Connected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connected
}

// SetConnected sets the logical client link. It never touches the sensors.
func (t *Amb) SetConnected(connected bool) {
	t.mu.Lock()
	t.connected = connected
	t.mu.Unlock()
}

// AveragePeriod returns the ObservingConditions averaging window, in hours.
func (t *Amb) AveragePeriod() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.averagePeriod
}

// SetAveragePeriod stores the ObservingConditions averaging window, in hours.
func (t *Amb) SetAveragePeriod(hours float64) {
	t.mu.Lock()
	t.averagePeriod = hours
	t.mu.Unlock()
}

func ambCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "AMB"))
}
