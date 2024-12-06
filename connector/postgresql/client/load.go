package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
)

func (c *Client) Migrate(ctx context.Context) {}

func (c *Client) Delete(ctx context.Context) {}

func (c *Client) Create(ctx context.Context, record chan<- arrow.Record) error {
	return nil
}
