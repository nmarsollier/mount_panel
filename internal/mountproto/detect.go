package mountproto

import (
	"context"
	"errors"
	"time"

	"github.com/nmarsollier/mount_panel/internal/proto"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// handshakeTimeout is how long a candidate port gets to answer the
// capabilities handshake during detection.
const handshakeTimeout = 1500 * time.Millisecond

var errNoMount = errors.New("no BFT mount found on any serial port")

// Detect probes the serial ports and returns a client connected to the BFT
// mount. Each candidate is claimed atomically before it is probed, so a random
// serial device is never mistaken for the mount and the accessory detector is
// never blocked by this one.
func Detect() (*Client, error) {
	proto.DetectMu.Lock()
	defer proto.DetectMu.Unlock()

	ports, err := proto.CandidatePorts()
	if err != nil {
		return nil, err
	}

	for _, name := range ports {
		if !proto.ClaimPort(name) {
			logx.Info(detectCtx(), "skip "+name+" (owned by another device)")
			continue
		}

		p, err := proto.OpenPort(name)
		if err != nil {
			proto.ReleasePort(name)
			continue
		}

		base := proto.NewClient(p, name, handshakeTimeout)
		c := &Client{Client: base}
		caps, err := c.Capabilities()
		if err == nil && isMount(caps) {
			base.SetTimeout(proto.DefaultTimeout)
			return c, nil
		}

		if err != nil {
			logx.Info(detectCtx(), "probe "+name+": "+err.Error())
		} else {
			logx.Info(detectCtx(), "probe "+name+": not a mount (device_type="+
				caps.DeviceType+")")
		}
		base.Close()
	}
	return nil, errNoMount
}

// isMount reports whether a capabilities response identifies the mount, by the
// shared device_type field.
func isMount(caps *CapabilitiesResponse) bool {
	return caps != nil && caps.DeviceType == "mount"
}

func detectCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "MountDetect"))
}
