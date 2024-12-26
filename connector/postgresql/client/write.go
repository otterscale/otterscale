package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
)

func (c *Client) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	// check namespace
	if c.opts.namespace == "" {
		return errors.New("namespace is empty")
	}
	// get current schema
	tables, err := c.GetTables(ctx, c.opts.namespace)
	if err != nil {
		return err
	}
	// create transaction
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// start writing
	if err := c.write(ctx, tx, tables, msgs); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			slog.Error("failed to rollback")
		}
		return err
	}
	// commit transaction
	return tx.Commit(ctx)
}

func (c *Client) write(ctx context.Context, tx pgx.Tx, tables []*arrow.Schema, msgs <-chan *pb.Message) error {
	for msg := range msgs {
		rec, err := openhdc.AppendBuiltinFieldsToRecord(msg)
		if err != nil {
			return err
		}
		switch msg.GetKind() {
		case property.MessageKind_migrate:
			if err := migrate(ctx, tx, tables, rec.Schema()); err != nil {
				return err
			}
		case property.MessageKind_insert:
			if err := insert(ctx, tx, rec, c.Codec); err != nil {
				return err
			}
		case property.MessageKind_upsert_update:
			if err := upsert(ctx, tx, rec, c.Codec, true); err != nil {
				return err
			}
		case property.MessageKind_upsert_nothing:
			if err := upsert(ctx, tx, rec, c.Codec, false); err != nil {
				return err
			}
		case property.MessageKind_delete_stale:
			if err := delete(ctx, tx, rec, msg.GetSyncedAt().AsTime()); err != nil {
				return err
			}
		case property.MessageKind_delete_all:
			if err := truncate(ctx, tx, rec.Schema()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("not supported kind %v", msg.GetKind())
		}
	}
	return nil
}
