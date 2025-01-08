package client

import (
	"context"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	return openhdc.ErrNotSupported
}
