package client

import (
	"context"
	"database/sql"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/connectors/oracle/client/or"
	"github.com/openhdc/openhdc/connectors/oracle/orarrow"
	"github.com/openhdc/openhdc/metadata"
)

type Tables []*arrow.Schema

func (ts Tables) Get(table string) (*arrow.Schema, bool) {
	for _, t := range ts {
		cur, _ := metadata.GetTableName(t)
		if cur == table {
			return t, true
		}
	}
	return nil, false
}

func newTables(ctx context.Context, pool *sql.DB) (Tables, error) {
	classes, err := or.Classes(ctx, pool)
	if err != nil {
		return nil, err
	}
	tab := Tables{}
	for _, class := range classes {
		atts, err := or.Attributes(ctx, pool, class.RelName)
		if err != nil {
			return nil, err
		}
		fs := []arrow.Field{}
		for _, att := range atts {
			fs = append(fs, arrow.Field{
				Name:     att.AttName,
				Type:     orarrow.ToForOID(att.AttTypeID),
				Nullable: !att.AttNotNull,
				Metadata: toFieldMetadata(att),
			})
		}
		md := toSchemaMetadata(class.RelName)
		sch := arrow.NewSchema(fs, md)
		tab = append(tab, sch)
	}
	return tab, nil
}
