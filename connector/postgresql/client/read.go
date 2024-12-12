package client

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/connector"
	"github.com/openhdc/openhdc/internal/metadata"
)

func (c *Client) Read(ctx context.Context, msg chan<- *pb.Message, opts connector.ReadOptions) error {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				slog.Error("failed to rollback")
			}
		}
	}()

	tables, err := c.GetTables(ctx, opts.Namespace)
	if err != nil {
		return err
	}

	for _, table := range tables {
		schema := table.Schema()
		if skip(schema, opts.Tables, opts.SkipTables) {
			continue
		}
		if err := c.migrate(schema, msg); err != nil {
			return err
		}
		if err := c.insert(ctx, schema, msg); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func skip(sch *arrow.Schema, tables, skipTables []string) bool {
	tableName, ok := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	if !ok {
		return true
	}
	if slices.Contains(skipTables, tableName) {
		return true
	}
	if len(tables) > 0 && !slices.Contains(tables, tableName) {
		return true
	}
	return false
}
