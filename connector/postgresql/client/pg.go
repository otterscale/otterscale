package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/internal/metadata"
)

type pgClass struct {
	RelName string `db:"relname"`
	OID     uint32 `db:"oid"`
}

type pgAttribute struct {
	AttName    string  `db:"attname"`
	AttTypeID  uint32  `db:"atttypid"`
	AttNotNull bool    `db:"attnotnull"`
	ConType    *rune   `db:"contype"`
	ConName    *string `db:"conname"`
}

func pgClasses(ctx context.Context, pool *pgxpool.Pool, namespace string) ([]*pgClass, error) {
	var sql = `select c.relname, c.oid
from pg_catalog.pg_namespace n
left join pg_catalog.pg_class c on n.oid = c.relnamespace 
where c.relkind = 'r' and n.nspname = $1
`
	rows, err := pool.Query(ctx, sql, namespace)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[pgClass])
}

func pgAttributes(ctx context.Context, pool *pgxpool.Pool, attrelid uint32) ([]*pgAttribute, error) {
	sql := `select a.attname, a.atttypid, a.attnotnull, c.contype, c.conname
from pg_catalog.pg_attribute a
left join pg_catalog.pg_constraint c on a.attrelid = c.conrelid and a.attnum = any(c.conkey)
where not a.attisdropped and a.attnum > 0 and a.attrelid = $1
order by a.attnum
`
	rows, err := pool.Query(ctx, sql, attrelid)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[pgAttribute])
}

func toFieldMetadata(att *pgAttribute, dataType string) arrow.Metadata {
	m := map[string]string{
		metadata.KeyFieldDataType: dataType,
	}
	if att.ConType != nil && att.ConName != nil {
		switch *att.ConType {
		case 'p':
			m[metadata.KeyFieldPrimaryKeyName] = *att.ConName
		case 'u':
			m[metadata.KeyFieldUniqueName] = *att.ConName
		}
	}
	return arrow.MetadataFrom(m)
}

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{
		metadata.KeySchemaTableName: tableName,
	}
	md := arrow.MetadataFrom(m)
	return &md
}
