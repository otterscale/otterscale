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
	"github.com/openhdc/openhdc/api/workload/v1"
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
		if skip(sch, opts.Keys, opts.SkipKeys) {
			continue
		}
		if err := c.read(ctx, tx, sch, msg, opts); err != nil {
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

func cursorsToWhere(syncMode property.SyncMode, cs []*workload.Sync_Option_Cursor) string {
	supports := []property.SyncMode{property.SyncMode_incremental_append, property.SyncMode_incremental_append_dedupe}
	if len(cs) == 0 || !slices.Contains(supports, syncMode) {
		return ""
	}
	ws := []string{}
	for _, c := range cs {
		ws = append(ws, fmt.Sprintf("%s > '%s'", sanitize(c.GetField()), c.GetValue()))
	}
	b := strings.Builder{}
	b.WriteString(" where ")
	b.WriteString(strings.Join(ws, " and "))
	return b.String()
}

func (c *Client) read(ctx context.Context, tx pgx.Tx, sch *arrow.Schema, msg chan<- *pb.Message, opts openhdc.ReadOptions) error {
	// timestamp
	syncedAt := time.Now().UTC().Truncate(time.Second)

	// record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// migration
	new, err := pb.NewMessage(property.MessageKind_migrate, builder.NewRecord(), c.opts.name, syncedAt)
	if err != nil {
		return err
	}
	msg <- new

	// get table
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	// sync mode
	syncMode := workload.GetSyncMode(opts.Options, tableName)
	cursors := workload.GetSyncCursors(opts.Options, tableName)

	// message kind
	kind, err := getMessageKind(sch, syncMode, cursors)
	if err != nil {
		return err
	}

	// truncate
	if deleteAll(sch, syncMode) {
		new, err := pb.NewMessage(property.MessageKind_delete_all, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	// get data
	columnNames := []string{}
	for _, field := range sch.Fields() {
		columnNames = append(columnNames, sanitize(field.Name))
	}

	// query
	b := strings.Builder{}
	b.WriteString("select ")
	b.WriteString(strings.Join(columnNames, ","))
	b.WriteString(" from ")
	b.WriteString(sanitize(tableName))
	b.WriteString(cursorsToWhere(syncMode, cursors))
	b.WriteString(" limit 1 ")
	rows, err := tx.Query(ctx, b.String())
	if err != nil {
		return err
	}
	defer rows.Close()

	// start
	var count int64
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

		count++
		if count > opts.BatchSize {
			new, err := pb.NewMessage(kind, builder.NewRecord(), c.opts.name, syncedAt)
			if err != nil {
				return err
			}
			msg <- new
			count = 0
		}
	}

	if count > 0 {
		new, err := pb.NewMessage(kind, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	// delete not exists
	if deleteStale(sch, syncMode) {
		new, err := pb.NewMessage(property.MessageKind_delete_stale, builder.NewRecord(), c.opts.name, syncedAt)
		if err != nil {
			return err
		}
		msg <- new
	}

	return rows.Err()
}

func getMessageKind(sch *arrow.Schema, syncMode property.SyncMode, cursors []*workload.Sync_Option_Cursor) (property.MessageKind, error) {
	hasPrimaryKey := metadata.HasPrimaryKey(sch)
	hasCursor := len(cursors) > 0
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
