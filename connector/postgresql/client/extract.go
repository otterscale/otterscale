package client

import (
	"context"
	"fmt"
	"reflect"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

var schemaSQL = ``

func (c *Client) GetTables(ctx context.Context) ([]arrow.Table, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tables := []arrow.Table{}

	return tables, nil
}

func (c *Client) List(ctx context.Context, record chan<- arrow.Record) error {
	sql := ""
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	// FIXME: STUB
	schema := &arrow.Schema{}
	codec := NewCodec()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return err
		}

		b := array.NewRecordBuilder(memory.DefaultAllocator, schema)
		defer b.Release()

		for idx, val := range vals {
			if err := codec.Append(b.Field(idx), val); err != nil {
				fmt.Println(reflect.TypeOf(b.Field(idx)), reflect.TypeOf(val))
				return err
			}
		}

		rec := b.NewRecord()
		defer rec.Release()

		record <- rec
	}

	return rows.Err()
}
