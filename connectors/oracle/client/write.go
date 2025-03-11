package client

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	return nil
}
