package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/connector"
)

func (c *Client) Read(ctx context.Context, rec chan<- arrow.Record, opts ...connector.ReadOption) error {
	return nil
}
