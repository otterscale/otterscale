package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
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
		if err := c.read(ctx, schema, msg); err != nil {
			slog.Error(err.Error())
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

func sanitize(str string) string {
	return pgx.Identifier{str}.Sanitize()
}

func (c *Client) read(ctx context.Context, sch *arrow.Schema, msg chan<- *pb.Message) error {
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	new, err := pb.NewMessage(pb.Kind_KIND_MIGRATE, builder.NewRecord())
	if err != nil {
		return err
	}
	msg <- new

	tableName, _ := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	columnNames := []string{}
	for _, field := range sch.Fields() {
		columnNames = append(columnNames, sanitize(field.Name))
	}

	sql := fmt.Sprintf("select %s from %s limit 1", strings.Join(columnNames, ","), sanitize(tableName))
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return err
		}

		for idx, val := range vals {
			if err := c.Append(builder.Field(idx), val); err != nil {
				slog.Error("invalid append", "type of field", reflect.TypeOf(builder.Field(idx)), "type of value", reflect.TypeOf(val))
				return err
			}
		}

		// TODO: BATCH
		new, err := pb.NewMessage(pb.Kind_KIND_INSERT, builder.NewRecord())
		if err != nil {
			return err
		}
		msg <- new
	}

	return rows.Err()
}
