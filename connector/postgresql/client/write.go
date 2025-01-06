package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/connector/postgresql/client/pg"
)

func (c *Client) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	// check namespace
	if c.opts.namespace == "" {
		return errors.New("namespace is empty")
	}
	// get current schema
	tables, err := newTables(ctx, c.pool, c.opts.namespace)
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

func (c *Client) write(ctx context.Context, tx pgx.Tx, tables Tables, msgs <-chan *pb.Message) error {
	for msg := range msgs {
		rec, err := openhdc.AppendBuiltinFieldsToRecord(msg)
		if err != nil {
			return err
		}
		h, err := pg.NewHelper(rec.Schema(), c.Codec)
		if err != nil {
			return err
		}
		switch msg.GetKind() {
		case property.MessageKind_migrate:
			if err := migrate(ctx, tx, tables, h); err != nil {
				return err
			}
		case property.MessageKind_insert:
			if err := h.Insert(ctx, tx, rec.NumRows(), rec.Columns()); err != nil {
				return err
			}
		case property.MessageKind_upsert_update:
			if err := h.Upsert(ctx, tx, rec.NumRows(), rec.Columns(), true); err != nil {
				return err
			}
		case property.MessageKind_upsert_nothing:
			if err := h.Upsert(ctx, tx, rec.NumRows(), rec.Columns(), false); err != nil {
				return err
			}
		case property.MessageKind_delete_stale:
			if err := h.Delete(ctx, tx, msg.GetSyncedAt().AsTime()); err != nil {
				return err
			}
		case property.MessageKind_delete_all:
			if err := h.Truncate(ctx, tx); err != nil {
				return err
			}
		default:
			return fmt.Errorf("not supported kind %v", msg.GetKind())
		}
	}
	return nil
}

func migrate(ctx context.Context, tx pgx.Tx, tables Tables, h *pg.Helper) error {
	cur, ok := tables.Get(h.TableName())
	if !ok {
		return h.CreateTable(ctx, tx)
	}
	add, del, ok := openhdc.CompareSchemata(cur, h.Schema())
	if !ok {
		return h.RenewTable(ctx, tx)
	}
	if len(add) > 0 || len(del) > 0 {
		if err := h.AlterTable(ctx, tx, add, del); err != nil {
			slog.Error(err.Error())
			return h.RenewTable(ctx, tx)
		}
	}
	return nil
}
