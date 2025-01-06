package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Class struct {
	RelName string `db:"relname"`
	OID     uint32 `db:"oid"`
}

func Classes(ctx context.Context, pool *pgxpool.Pool, namespace string) ([]*Class, error) {
	sql := `select c.relname, c.oid
from pg_catalog.pg_namespace n
left join pg_catalog.pg_class c on n.oid = c.relnamespace 
where c.relkind = 'r' and n.nspname = $1
`
	rows, err := pool.Query(ctx, sql, namespace)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Class])
}

type Attribute struct {
	AttNum     string  `db:"attnum"`
	AttName    string  `db:"attname"`
	AttTypeID  uint32  `db:"atttypid"`
	AttNotNull bool    `db:"attnotnull"`
	ConTypes   *string `db:"contypes"`
}

func Attributes(ctx context.Context, pool *pgxpool.Pool, attrelid uint32) ([]*Attribute, error) {
	sql := `select a.attnum, a.attname, a.atttypid, a.attnotnull, string_agg(c.contype, ',') contypes
from pg_catalog.pg_attribute a
left join pg_catalog.pg_constraint c on a.attrelid = c.conrelid and a.attnum = any(c.conkey)
where not a.attisdropped and a.attnum > 0 and a.attrelid = $1
group by a.attnum, a.attname, a.atttypid, a.attnotnull
order by a.attnum
`
	rows, err := pool.Query(ctx, sql, attrelid)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Attribute])
}
