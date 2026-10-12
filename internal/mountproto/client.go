package mountproto

import (
	"context"
	"encoding/json"

	"github.com/nmarsollier/nrm_ha_go/internal/proto"
)

// Client talks to the mount over the USB CDC-ACM protocol. It wraps the shared
// transport with the mount message contract: flat {type, id, ...fields}
// requests and {ok, error:string} rejections.
type Client struct {
	*proto.Client
}

// Open connects to the given serial port.
func Open(name string) (*Client, error) {
	c, err := proto.Open(name)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// request performs a round trip and unmarshals the response into out, checking
// the shared ok/error envelope first.
func (c *Client) request(ctx context.Context, reqType string, payload any, out any) error {
	id := c.NextID()
	body, err := buildRequest(reqType, id, payload)
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
		return &Error{Message: env.Error}
	}

	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// buildRequest flattens a family payload into the common request envelope.
func buildRequest(reqType string, id uint32, payload any) ([]byte, error) {
	m := map[string]any{"type": reqType, "id": id}

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		var fields map[string]any
		if err := json.Unmarshal(b, &fields); err != nil {
			return nil, err
		}
		for k, v := range fields {
			m[k] = v
		}
	}

	return json.Marshal(m)
}
