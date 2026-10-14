package eafproto

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

var errNoAccessories = errors.New("no NRM accessory controller found on any serial port")

// Detect probes the serial ports and returns a client connected to the
// accessory controller (device_type="accessories"). Each candidate is claimed
// atomically before it is probed, so a random device is never mistaken for the
// accessory and the mount detector is never blocked by this one.
func Detect() (*Client, error) {
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
		c := &Client{Client: base, session: newSession()}
		caps, err := c.Capabilities()
		if err == nil && isAccessories(caps) {
			c.bootID = caps.BootID
			base.SetTimeout(proto.DefaultTimeout)
			return c, nil
		}

		if err != nil {
			logx.Info(detectCtx(), "probe "+name+": "+err.Error())
		} else {
			logx.Info(detectCtx(), "probe "+name+": not an accessory (device_type="+
				caps.DeviceType+")")
		}
		base.Close()
	}
	return nil, errNoAccessories
}

// isAccessories reports whether a capabilities response identifies the
// accessory controller, by the shared device_type field.
func isAccessories(caps *CapabilitiesResponse) bool {
	return caps != nil && caps.DeviceType == "accessories"
}

func detectCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "EAFDetect"))
}
