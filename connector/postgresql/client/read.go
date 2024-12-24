package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/metadata"
)

func (c *Client) Read(ctx context.Context, msg chan<- *pb.Message, opts openhdc.ReadOptions) error {
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

	if c.opts.namespace == "" {
		return errors.New("namespace is empty")
	}

	schs, err := c.GetTables(ctx, c.opts.namespace)
	if err != nil {
		return err
	}

	for _, sch := range schs {
		if skip(sch, opts.Tables, opts.SkipTables) {
			continue
		}
		if err := c.read(ctx, tx, sch, msg); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func skip(sch *arrow.Schema, tables, skipTables []string) bool {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
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

func (c *Client) read(ctx context.Context, tx pgx.Tx, sch *arrow.Schema, msg chan<- *pb.Message) error {
	// timestamp
	syncedAt := time.Now().UTC().Truncate(time.Second)

	// message kind
	kind, err := getMessageKind(sch, c.opts.cursor, c.opts.syncMode)
	if err != nil {
		return err
	}

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// migration
	new, err := pb.NewMessage(property.MessageKind_migrate, builder.NewRecord(), c.opts.name, syncedAt)
	if err != nil {
		return err
	}
	msg <- new

	// truncate
	if deleteAll(sch, c.opts.syncMode) {
		new, err := pb.NewMessage(property.MessageKind_delete_all, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	// get table
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	// get columns
	columnNames := []string{}
	for _, field := range sch.Fields() {
		columnNames = append(columnNames, sanitize(field.Name))
	}

	// get data
	sql := fmt.Sprintf("select %s from %s limit 1", strings.Join(columnNames, ","), sanitize(tableName))
	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	// start
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return err
		}

		for idx, val := range vals {
			if err := c.Encode(builder.Field(idx), val); err != nil {
				slog.Error("invalid append", "type of field", reflect.TypeOf(builder.Field(idx)), "type of value", reflect.TypeOf(val))
				return err
			}
		}

		new, err := pb.NewMessage(kind, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	// delete not exists
	if deleteStale(sch, c.opts.syncMode) {
		new, err := pb.NewMessage(property.MessageKind_delete_stale, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	return rows.Err()
}

func getMessageKind(sch *arrow.Schema, cursor string, syncMode property.SyncMode) (property.MessageKind, error) {
	hasPrimaryKey := metadata.HasPrimaryKey(sch)
	hasCursor := cursor != ""
	switch syncMode {
	case property.SyncMode_full_overwrite:
		// upsert & delete
		if hasPrimaryKey {
			return property.MessageKind_upsert_update, nil
		}
		// truncate & insert
		return property.MessageKind_insert, nil

	case property.SyncMode_full_append:
		// do nothing
		if hasPrimaryKey {
			return property.MessageKind_upsert_nothing, nil
		}
		// insert
		return property.MessageKind_insert, nil

	case property.SyncMode_incremental_append:
		if !hasCursor {
			return 0, errors.New("cursor is empty")
		}
		if hasPrimaryKey {
			return 0, errors.New("primary key is exists, use `incremental_append_dedupe` instead")
		}
		// insert
		return property.MessageKind_insert, nil

	case property.SyncMode_incremental_append_dedupe:
		if !hasCursor {
			return 0, errors.New("cursor is empty")
		}
		if !hasPrimaryKey {
			return 0, errors.New("primary key is empty, use `incremental_append` instead")
		}
		// upsert_update
		return property.MessageKind_upsert_update, nil

	default:
		return property.MessageKind_message_kind_unspecified, nil
	}
}

func deleteAll(sch *arrow.Schema, syncMode property.SyncMode) bool {
	return syncMode == property.SyncMode_full_overwrite && !metadata.HasPrimaryKey(sch)
}

func deleteStale(sch *arrow.Schema, syncMode property.SyncMode) bool {
	return syncMode == property.SyncMode_full_overwrite && metadata.HasPrimaryKey(sch)
}
