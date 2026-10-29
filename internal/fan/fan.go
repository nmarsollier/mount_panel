// Package fan models the primary-mirror fan controller that hangs off the BFT
// mount over its own USB CDC-ACM link. It mirrors the eaf package: a singleton
// that keeps the serial link up and translates high-level operations into
// eafproto calls. No state cache: every read fetches a fresh snapshot.
package fan

import (
	"context"
	"errors"
	"sync"

	"github.com/nmarsollier/mount_panel/internal/eafproto"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// errNotConnected is returned when an operation is requested but the serial
// link to the fan controller is down.
var errNotConnected = errors.New("fan controller not connected")

// Fan is the model of the fan controller: it manages the serial connection and
// exposes the fan and mirror-temperature operations to the UI.
type Fan struct {
	mu sync.RWMutex

	conn       *eafproto.Client // nil = disconnected
	serialPort string
	connecting bool

	connected bool // logical client link (Alpaca Connected)
}

var (
	once sync.Once
	inst *Fan
)

// Get returns the singleton fan model. It is non-nil after Init.
func Get() *Fan {
	return inst
}

// Init starts the fan model. It makes a single startup connection attempt;
// further attempts are manual, from the UI, via Connect.
func Init(serialPort string) {
	once.Do(func() {
		inst = &Fan{serialPort: serialPort}
		go inst.Connect()
	})
}

// Connect makes a single attempt to open and verify the fan controller. It is a
// no-op when the link is already up or a connection attempt is in flight, and
// returns the error when the attempt fails.
func (f *Fan) Connect() error {
	f.mu.Lock()
	if f.conn != nil || f.connecting {
		f.mu.Unlock()
		return nil
	}
	f.connecting = true
	f.mu.Unlock()
	defer func() {
		f.mu.Lock()
		f.connecting = false
		f.mu.Unlock()
	}()

	c, err := f.tryConnect()
	if err != nil {
		logx.Warn(fanCtx(), "connect failed: "+err.Error())
		return err
	}
	go f.watch(c)
	return nil
}

// Disconnect closes the serial link to the fan controller. It is a no-op when
// the link is already down.
func (f *Fan) Disconnect() {
	f.mu.RLock()
	c := f.conn
	f.mu.RUnlock()
	if c == nil {
		return
	}
	f.disconnect(c)
}

// watch tears the link down when the reader notices it has dropped. The link is
// not reconnected automatically: the user re-connects from the UI.
func (f *Fan) watch(c *eafproto.Client) {
	<-c.Done()
	reason := ""
	if err := c.Err(); err != nil {
		reason = ": " + err.Error()
	}
	logx.Warn(fanCtx(), "fan disconnected"+reason)
	f.disconnect(c)
}

// tryConnect opens and verifies the fan controller.
func (f *Fan) tryConnect() (*eafproto.Client, error) {
	var (
		c   *eafproto.Client
		err error
	)
	if f.serialPort != "" {
		c, err = eafproto.Open(f.serialPort)
	} else {
		c, err = eafproto.Detect("fan")
	}
	if err != nil {
		return nil, err
	}

	f.mu.Lock()
	f.conn = c
	f.mu.Unlock()

	logx.Info(fanCtx(), "connected on "+c.PortName()+" boot="+c.BootID())
	return c, nil
}

// disconnect closes the client and clears the connection.
func (f *Fan) disconnect(c *eafproto.Client) {
	c.Close()
	f.mu.Lock()
	if f.conn == c {
		f.conn = nil
	}
	f.mu.Unlock()
}

// getConn returns the current serial connection, or an error when disconnected.
func (f *Fan) getConn() (*eafproto.Client, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.conn == nil {
		return nil, errNotConnected
	}
	return f.conn, nil
}

// IsConnected reports whether the serial link to the fan controller is up.
func (f *Fan) IsConnected() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.conn != nil
}

// Connected reports the logical client link (Alpaca Connected).
func (f *Fan) Connected() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.connected
}

// SetConnected sets the logical client link. It never touches the fan.
func (f *Fan) SetConnected(connected bool) {
	f.mu.Lock()
	f.connected = connected
	f.mu.Unlock()
}

func fanCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "FAN"))
}
