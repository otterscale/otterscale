package client

import (
	"context"
	"fmt"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/connector"
)

func (c *Client) Write(ctx context.Context, msg <-chan *pb.Message, opts connector.WriteOptions) error {
	// TODO: FAKE
	namespace := "public"

	tables, err := c.GetTables(ctx, namespace)
	if err != nil {
		return err
	}

	for {
		select {
		case msg, ok := <-msg:
			if !ok {
				continue // ?
			}
			rec, err := pb.ToArrowRecord(msg.Record)
			if err != nil {
				return err
			}
			switch msg.Kind {
			case pb.Kind_KIND_MIGRATE:
				if err := c.migrate(ctx, tables, rec.Schema()); err != nil {
					return err
				}

			case pb.Kind_KIND_INSERT:

			case pb.Kind_KIND_DELETE:

			default:
				return fmt.Errorf("not supported kind %v", msg.Kind)
			}
		}
	}
}
