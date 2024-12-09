package client

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) GetTable(ctx context.Context, table string) (arrow.Table, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	t, err := conn.Conn().LoadType(ctx, table)
	if err != nil {
		return nil, err
	}

	cc, ok := t.Codec.(*pgtype.CompositeCodec)
	if !ok {
		return nil, fmt.Errorf("invalid pgtype %T", t.Codec)
	}

	columns := []arrow.Column{}
	fields := []arrow.Field{}

	for _, f := range cc.Fields {
		arr, err := c.Encode(f.Type.Codec, nil)
		if err != nil {
			return nil, err
		}
		field := arrow.Field{
			Name: f.Name,
			Type: arr.DataType(),
			// Nullable: !attribute.IsNotNull,
			// Metadata: arrow.NewMetadata(
			// 	[]string{
			// 		MetadataFieldKeyBaseColumnName,
			// 		MetadataFieldKeyDataType,
			// 		MetadataFieldKeyPKName,
			// 		MetadataFieldKeyUniqueName,
			// 	},
			// 	[]string{
			// 		attribute.ColumnName,
			// 		attribute.DataType,
			// 		attribute.PKName,
			// 		attribute.UniqueName,
			// 	}),
		}
		fields = append(fields, field)
		chunked := arrow.NewChunked(arr.DataType(), nil)
		column := arrow.NewColumn(field, chunked)
		columns = append(columns, *column)
	}

	schema := arrow.NewSchema(fields, nil)
	return array.NewTable(schema, columns, 0), nil
}
