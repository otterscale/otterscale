package client

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/openhdc/openhdc/internal/connector"
)

var schemaSQL = ``

func (c *Client) Read(ctx context.Context, rec chan<- arrow.Record, opts ...connector.ReadOption) error {
	sql := ""
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	// FIXME: STUB
	t, err := c.GetTable(ctx, "test_case")
	if err != nil {
		return err
	}
	schema := t.Schema()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return err
		}

		b := array.NewRecordBuilder(memory.DefaultAllocator, schema)
		defer b.Release()

		for idx, val := range vals {
			if err := c.Append(b.Field(idx), val); err != nil {
				slog.Error("invalid append", "type of field", reflect.TypeOf(b.Field(idx)), "type of value", reflect.TypeOf(val))
				return err
			}
		}

		r := b.NewRecord()
		defer r.Release()

		rec <- r
	}

	return rows.Err()
}
