package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/connector"
)

func (c *Client) Write(ctx context.Context, rec chan<- arrow.Record, opts ...connector.WriteOption) error {
	return nil
}
