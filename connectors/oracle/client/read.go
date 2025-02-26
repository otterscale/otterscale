package client

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"reflect"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/connectors/oracle/client/or"
)

func (c *Client) Read(ctx context.Context, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// check
	if c.opts.namespace == "" {
		return errors.New("namespace is empty")
	}

	// new transaction
	tx, err := c.pool.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				slog.Error("failed to rollback")
			}
		}
	}()

	// new tables
	tables, err := newTables(ctx, c.pool)
	if err != nil {
		return err
	}

	// start
	for _, sch := range tables {
		if skip(sch, rdr.Keys(), rdr.SkipKeys()) {
			continue
		}
		if err := c.read(ctx, tx, sch, msgs, rdr); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) read(ctx context.Context, tx *sql.Tx, sch *arrow.Schema, msgs chan<- *pb.Message, rdr *openhdc.Reader) error {
	// record builder
	b := array.NewRecordBuilder(memory.DefaultAllocator, sch)

	// migration
	if err := rdr.Send(pb.Migrate, msgs, b.NewRecord()); err != nil {
		return err
	}

	// new helper
	h, err := or.NewHelper(sch, c.Codec)
	if err != nil {
		return err
	}

	// sync mode
	mode := rdr.SyncMode(h.TableName())
	curs := rdr.SyncCursors(h.TableName())

	// message kind
	kind, err := toMessageKind(sch, mode, curs)
	if err != nil {
		return err
	}

	// truncate
	if deleteAll(sch, mode) {
		if err := rdr.Send(pb.DeleteAll, msgs, b.NewRecord()); err != nil {
			return err
		}
	}

	// query
	rows, err := h.Select(ctx, tx, mode, curs)
	if err != nil {
		return err
	}
	defer rows.Close()

	// columns
	cols, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	// start
	var count int64
	for rows.Next() {
		// query
		ptrs := make([]any, len(cols))
		vals := make([]sql.NullString, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		// encode
		for idx, val := range vals {
			if err := c.Encode(b.Field(idx), val.String); err != nil {
				slog.Error("invalid append", "type of field", reflect.TypeOf(b.Field(idx)), "type of value", reflect.TypeOf(val.String))
				return err
			}
		}

		// batch
		count++
		if count > rdr.BatchSize() {
			if err := rdr.Send(kind, msgs, b.NewRecord()); err != nil {
				return err
			}
			count = 0
		}
	}

	// remain
	if count > 0 {
		if err := rdr.Send(kind, msgs, b.NewRecord()); err != nil {
			return err
		}
	}

	// delete not exists
	if deleteStale(sch, mode) {
		if err := rdr.Send(pb.DeleteStale, msgs, b.NewRecord()); err != nil {
			return err
		}
	}

	return rows.Err()
}
