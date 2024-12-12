package client

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) GetTables(ctx context.Context, namespace string) ([]arrow.Table, error) {
	m := pgtype.NewMap()
	classes, err := pgClasses(ctx, c.pool, namespace)
	if err != nil {
		return nil, err
	}
	tables := []arrow.Table{}
	for _, class := range classes {
		atts, err := pgAttributes(ctx, c.pool, class.OID)
		if err != nil {
			return nil, err
		}
		fields := []arrow.Field{}
		for _, att := range atts {
			typ, ok := m.TypeForOID(att.AttTypeID)
			if !ok {
				return nil, fmt.Errorf("invalid pgtype %d", att.AttTypeID)
			}
			arr, err := c.Encode(typ.Codec, nil)
			if err != nil {
				return nil, err
			}
			fields = append(fields, arrow.Field{
				Name:     att.AttName,
				Type:     arr.DataType(),
				Nullable: !att.AttNotNull,
				Metadata: toFieldMetadata(att, typ.Name),
			})
		}
		columns := []arrow.Column{}
		for _, field := range fields {
			columns = append(columns, *arrow.NewColumn(field, arrow.NewChunked(field.Type, nil)))
		}
		schema := arrow.NewSchema(fields, toSchemaMetadata(class.RelName))
		tables = append(tables, array.NewTable(schema, columns, -1))
	}
	return tables, nil
}
