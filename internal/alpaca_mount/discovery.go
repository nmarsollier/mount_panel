package alpaca_mount

import (
	"context"
	"encoding/json"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"net"
	"strings"

	"github.com/nmarsollier/mount_panel/internal/mount"
	"github.com/nmarsollier/mount_panel/internal/tools/env"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
)

// discoveryPort is the fixed ASCOM Alpaca Discovery Protocol port.
const discoveryPort = 32227

// startDiscovery launches the Alpaca Discovery listener (UDP 32227) in the
// background. A client sends the magic string "alpacadiscovery1" and receives
// a JSON payload carrying the Alpaca port.
func startDiscovery() {
	go discoveryLoop()
}

func discoveryLoop() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: discoveryPort})
	if err != nil {
		logx.Warn(discoveryCtx(), "UDP discovery unavailable: "+err.Error())
		return
	}
	defer conn.Close()

	buf := make([]byte, 128)
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		if !strings.Contains(string(buf[:n]), "alpacadiscovery1") {
			continue
		}
		// Suppress the response when the mount is in the error state, matching
		// the firmware's behaviour.
		if isErr, err := mount.Get().IsError(); err == nil && isErr {
			continue
		}
		_, _ = conn.WriteToUDP(discoveryResponse(), remote)
	}
}

func discoveryResponse() []byte {
	resp := struct {
		AlpacaPort       int    `json:"AlpacaPort"`
		ServerName       string `json:"ServerName"`
		Version          string `json:"Version"`
		InterfaceVersion int    `json:"InterfaceVersion"`
	}{
		AlpacaPort:       env.Get().Port,
		ServerName:       server.ServerName,
		Version:          "v1",
		InterfaceVersion: server.InterfaceVersion,
	}
	b, _ := json.Marshal(resp)
	return b
}

func discoveryCtx() context.Context {
	return logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Discovery"))
}
