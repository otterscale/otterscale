package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
)

func (c *Client) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	if c.opts.syncMode == property.SyncMode_sync_mode_unspecified {
		slog.Warn("sync mode is unspecified, use full_overwrite")
		c.opts.syncMode = property.SyncMode_full_overwrite
	}
	if c.opts.namespace == "" {
		return errors.New("namespace is empty")
	}
	tables, err := c.GetTables(ctx, c.opts.namespace)
	if err != nil {
		return err
	}
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	if err := c.write(ctx, tx, tables, msgs); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			slog.Error("failed to rollback")
		}
		return err
	}
	return tx.Commit(ctx)
}

func (c *Client) write(ctx context.Context, tx pgx.Tx, tables []*arrow.Schema, msgs <-chan *pb.Message) error {
	// TODO: BATCH
	for msg := range msgs {
		rec, err := pb.ToArrowRecord(msg.GetRecord())
		if err != nil {
			return err
		}
		switch msg.GetKind() {
		case property.MessageKind_migrate:
			if err := migrate(ctx, tx, tables, rec.Schema()); err != nil {
				return err
			}
		case property.MessageKind_insert:
			if err := insert(ctx, tx, c.Codec, rec); err != nil {
				return err
			}
		case property.MessageKind_delete: // TODO: AFTER WRITE MODE
		default:
			return fmt.Errorf("not supported kind %v", msg.GetKind())
		}
	}
	return nil
}
