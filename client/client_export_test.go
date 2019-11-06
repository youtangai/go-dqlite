package client

import (
	"github.com/CanonicalLtd/go-dqlite/internal/protocol"
)

func (c *Client) Protocol() *protocol.Protocol {
	return c.protocol
}
