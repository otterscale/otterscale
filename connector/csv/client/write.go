package client

import (
	"context"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Write(ctx context.Context, msg <-chan *pb.Message, opts openhdc.WriteOptions) error {
	return nil
}
