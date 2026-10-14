package eafproto

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/nmarsollier/mount_panel/internal/proto"
)

// protoVersion is the accessory protocol version.
const protoVersion = 2

// Client talks to the accessory controller (focuser, fan, environment) over the
// USB CDC-ACM protocol. It wraps the shared transport with the accessory
// contract: {proto, type, id, session, boot_id, target, action|op, params}
// requests and {ok, data | error{code,message,field}} responses.
type Client struct {
	*proto.Client

	session string
	bootID  string
}

// Open connects to the given serial port, without the role handshake.
func Open(name string) (*Client, error) {
	c, err := proto.Open(name)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c, session: newSession()}, nil
}

// newSession returns a fresh random session identifier for this Go instance.
func newSession() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "g-" + hex.EncodeToString(b)
}

// BootID returns the boot identifier reported at connect time.
func (c *Client) BootID() string { return c.bootID }

// request builds a request from the given fields and unmarshals a typed
// response, checking the shared ok/error envelope first.
func (c *Client) request(ctx context.Context, reqType, target, verb string, params map[string]any, out any) error {
	id := c.NextID()

	m := buildRequest(reqType, target, verb, c.session, c.bootID, params)
	m["id"] = id

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	raw, err := c.RoundTrip(ctx, id, body)
	if err != nil {
		return err
	}

	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	if !env.OK {
		return parseError(env.Error)
	}
	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// buildRequest assembles the accessory request envelope: {proto, type, target,
// action|op, params} plus the session and boot identity (id is added by the
// caller). It is a pure helper so the wire shape can be tested without a port.
func buildRequest(reqType, target, verb, session, bootID string, params map[string]any) map[string]any {
	m := map[string]any{
		"proto":   protoVersion,
		"type":    reqType,
		"session": session,
		"boot_id": bootID,
	}
	if target != "" {
		m["target"] = target
	}
	switch reqType {
	case "action":
		if verb != "" {
			m["action"] = verb
		}
	case "control", "config":
		if verb != "" {
			m["op"] = verb
		}
	}
	if params != nil {
		m["params"] = params
	}
	return m
}
