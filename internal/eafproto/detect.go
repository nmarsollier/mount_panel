package eafproto

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nmarsollier/mount_panel/internal/proto"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// handshakeTimeout is how long a candidate port gets to answer the
// capabilities handshake during detection.
const handshakeTimeout = 1500 * time.Millisecond

var errNoAccessories = errors.New("no accessory controller found on any serial port")

// Detect probes the serial ports and returns a client connected to the
// accessory controller whose capabilities report the given device_id ("eaf",
// "fan" or "temp"). Each candidate is claimed atomically before it is probed,
// so a random device is never mistaken for the accessory and the mount detector
// is never blocked by this one.
func Detect(deviceID string) (*Client, error) {
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
		c := &Client{Client: base, session: newSession()}
		caps, err := c.Capabilities()
		if err == nil && isDevice(caps, deviceID) {
			c.bootID = caps.BootID
			base.SetTimeout(proto.DefaultTimeout)
			return c, nil
		}

		if err != nil {
			logx.Info(detectCtx(), "probe "+name+": "+err.Error())
		} else {
			logx.Info(detectCtx(), "probe "+name+": not a "+deviceID+" (device_id="+
				caps.DeviceID+")")
		}
		base.Close()
	}
	return nil, fmt.Errorf("%w (device_id=%s)", errNoAccessories, deviceID)
}

// isDevice reports whether a capabilities response identifies the expected
// accessory controller, by the shared device_type field plus the device_id.
func isDevice(caps *CapabilitiesResponse, deviceID string) bool {
	return caps != nil && caps.DeviceType == "accessories" && caps.DeviceID == deviceID
}

func detectCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "EAFDetect"))
}
