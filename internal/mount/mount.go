package mount

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nmarsollier/nrm_ha_go/internal/proto"
	"github.com/nmarsollier/nrm_ha_go/internal/tools/logx"
)

// errNotConnected is returned when an operation is requested but the serial
// link to the mount is down.
var errNotConnected = errors.New("mount not connected")

// retryInterval is how long to wait before re-attempting a failed connection.
const retryInterval = 2 * time.Second

// Mount is the model of the telescope: it manages the serial connection and
// translates high-level operations into protocol calls. It does not cache the
// ESP state — every read fetches a fresh snapshot over the serial link.
type Mount struct {
	mu sync.RWMutex

	conn *proto.Client // nil = disconnected

	// App-owned transient state (reset on restart, not persisted).
	targetRa  float64
	targetDec float64

	selectedTrackingRate int // ASCOM DriveRates: 0 sidereal, 1 lunar, 2 solar
	connected            bool

	serialPort string
}

var (
	once sync.Once
	inst *Mount
)

// Get returns the singleton mount. It is non-nil after Init.
func Get() *Mount {
	return inst
}

// Init starts the mount model and its background loops (connection manager and
// time sync). It does not block: the mount connects whenever the ESP becomes
// available.
func Init(serialPort string) {
	once.Do(func() {
		inst = &Mount{serialPort: serialPort, selectedTrackingRate: 0}
		go inst.connectLoop()
		go inst.timeSyncLoop()
	})
}

// connectLoop keeps the serial link up: it connects, waits for the link to
// drop, then reconnects.
func (m *Mount) connectLoop() {
	for {
		c, err := m.tryConnect()
		if err != nil {
			logx.Warn(mountCtx(), "connect failed: "+err.Error())
			time.Sleep(retryInterval)
			continue
		}

		<-c.Done() // block until the link drops
		logx.Warn(mountCtx(), "mount disconnected")
		m.disconnect(c)
	}
}

// tryConnect opens and verifies the mount, then provides the UTC clock.
func (m *Mount) tryConnect() (*proto.Client, error) {
	var (
		c   *proto.Client
		err error
	)
	if m.serialPort != "" {
		c, err = proto.Open(m.serialPort)
	} else {
		c, err = proto.Detect()
	}
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.conn = c
	m.mu.Unlock()

	if err := m.provideTime(); err != nil {
		m.disconnect(c)
		return nil, err
	}
	return c, nil
}

// disconnect closes the client and clears the connection.
func (m *Mount) disconnect(c *proto.Client) {
	c.Close()
	m.mu.Lock()
	if m.conn == c {
		m.conn = nil
	}
	m.mu.Unlock()
}

// getConn returns the current serial connection, or an error when disconnected.
func (m *Mount) getConn() (*proto.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.conn == nil {
		return nil, errNotConnected
	}
	return m.conn, nil
}

// action sends an action request on the current connection.
func (m *Mount) action(req proto.ActionRequest) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	return c.Action(req)
}

// control sends a control request on the current connection.
func (m *Mount) control(req proto.ControlRequest) error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	return c.Control(req)
}

// IsConnected reports whether the serial link to the mount is up.
func (m *Mount) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conn != nil
}

// timeSyncLoop re-sends the UTC clock hourly to correct the ESP RTC drift.
func (m *Mount) timeSyncLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.provideTime(); err != nil && !errors.Is(err, errNotConnected) {
			logx.Warn(mountCtx(), "time re-sync failed: "+err.Error())
		}
	}
}

// fetchState returns a fresh state snapshot, or nil when unavailable.
func (m *Mount) fetchState() *proto.StateResponse {
	c, err := m.getConn()
	if err != nil {
		return nil
	}
	st, err := c.State()
	if err != nil {
		return nil
	}
	return st
}

// fetchConfig returns the current configuration, or nil when unavailable.
func (m *Mount) fetchConfig() *proto.ConfigResponse {
	c, err := m.getConn()
	if err != nil {
		return nil
	}
	cfg, err := c.ConfigGet()
	if err != nil {
		return nil
	}
	return cfg
}

// fetchCaps returns the mount capabilities, or nil when unavailable.
func (m *Mount) fetchCaps() *proto.CapabilitiesResponse {
	c, err := m.getConn()
	if err != nil {
		return nil
	}
	caps, err := c.Capabilities()
	if err != nil {
		return nil
	}
	return caps
}

// provideTime sends the current UTC clock so GOTO works now that the mount has
// no network of its own.
func (m *Mount) provideTime() error {
	c, err := m.getConn()
	if err != nil {
		return err
	}
	utc := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	_, err = c.ConfigSet(proto.ConfigRequest{Utc: &utc})
	return err
}

// IsError reports whether the mount is in the unrecoverable error state.
func (m *Mount) IsError() bool {
	s := m.fetchState()
	return s != nil && s.State == "error"
}

// Connected reports the logical client link (ASCOM Connected).
func (m *Mount) Connected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

// SetConnected sets the logical client link. It never moves the mount.
func (m *Mount) SetConnected(connected bool) {
	m.mu.Lock()
	m.connected = connected
	m.mu.Unlock()
}

func mountCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Mount"))
}
