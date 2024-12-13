package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/openhdc/openhdc/connector/postgresql/pgarrow"
)

func (c *Client) GetTables(ctx context.Context, namespace string) ([]arrow.Table, error) {
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
			fields = append(fields, arrow.Field{
				Name:     att.AttName,
				Type:     pgarrow.ToForOID(att.AttTypeID),
				Nullable: !att.AttNotNull,
				Metadata: toFieldMetadata(att),
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
