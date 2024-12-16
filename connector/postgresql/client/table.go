package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/connector/postgresql/pgarrow"
)

func (c *Client) GetTables(ctx context.Context, namespace string) ([]*arrow.Schema, error) {
	classes, err := pgClasses(ctx, c.pool, namespace)
	if err != nil {
		return nil, err
	}
	schs := []*arrow.Schema{}
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
		schs = append(schs, arrow.NewSchema(fields, toSchemaMetadata(class.RelName)))
	}
	return schs, nil
}
