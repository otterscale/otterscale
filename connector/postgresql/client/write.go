package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

func (c *Client) Write(ctx context.Context, msg <-chan *pb.Message, opts openhdc.WriteOptions) error {
	tables, err := c.GetTables(ctx, c.opts.namespace)
	if err != nil {
		return err
	}
	// TODO: BATCH
	for {
		msg, ok := <-msg
		if !ok {
			return errors.New("something wrong")
		}
		rec, err := pb.ToArrowRecord(msg.Record)
		if err != nil {
			return err
		}
		switch msg.Kind {
		case pb.Kind_KIND_MIGRATE:
			if err := c.migrate(ctx, tables, rec.Schema()); err != nil {
				slog.Error(err.Error())
				return err
			}
		case pb.Kind_KIND_INSERT:
			if err := c.insert(ctx, rec); err != nil {
				slog.Error(err.Error())
				return err
			}
		case pb.Kind_KIND_DELETE: // TODO: AFTER WRITE MODE
		default:
			return fmt.Errorf("not supported kind %v", msg.Kind)
		}
	}
}
