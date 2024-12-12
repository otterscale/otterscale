package client

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/jackc/pgx/v5"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/metadata"
)

func sanitize(str string) string {
	return pgx.Identifier{str}.Sanitize()
}

func (c *Client) insert(ctx context.Context, sch *arrow.Schema, msg chan<- *pb.Message) error {
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sch)

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
