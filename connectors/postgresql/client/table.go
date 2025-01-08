package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/connectors/postgresql/client/pg"
	"github.com/openhdc/openhdc/connectors/postgresql/pgarrow"
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

func newTables(ctx context.Context, pool *pgxpool.Pool, namespace string) (Tables, error) {
	classes, err := pg.Classes(ctx, pool, namespace)
	if err != nil {
		return nil, err
	}
	tab := Tables{}
	for _, class := range classes {
		atts, err := pg.Attributes(ctx, pool, class.OID)
		if err != nil {
			return nil, err
		}
		fs := []arrow.Field{}
		for _, att := range atts {
			fs = append(fs, arrow.Field{
				Name:     att.AttName,
				Type:     pgarrow.ToForOID(att.AttTypeID),
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
