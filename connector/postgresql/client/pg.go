package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/connector/postgresql/pgarrow"
	"github.com/openhdc/openhdc/metadata"
)

type pgClass struct {
	RelName string `db:"relname"`
	OID     uint32 `db:"oid"`
}

type pgAttribute struct {
	AttNum     string  `db:"attnum"`
	AttName    string  `db:"attname"`
	AttTypeID  uint32  `db:"atttypid"`
	AttNotNull bool    `db:"attnotnull"`
	ConTypes   *string `db:"contypes"`
}

func pgClasses(ctx context.Context, pool *pgxpool.Pool, namespace string) ([]*pgClass, error) {
	sql := `select c.relname, c.oid
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
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[pgAttribute])
}

func createTableIfNotExists(ctx context.Context, tx pgx.Tx, sch *arrow.Schema) error {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	pks := []string{}
	css := []string{}
	for _, field := range sch.Fields() {
		if metadata.IsPrimaryKey(&field) {
			pks = append(pks, sanitize(field.Name))
		}
		css = append(css, addColumnStatement(&field, false))
	}

	var b strings.Builder
	b.WriteString("create table if not exists ")
	b.WriteString(tableName)
	b.WriteString(" (")

	b.WriteString(strings.Join(css, ", "))

	if len(pks) > 0 {
		b.WriteString(", constraint ")
		b.WriteString(tableName)
		b.WriteString("_pkey")
		b.WriteString(" primary key (")
		b.WriteString(strings.Join(pks, ", "))
		b.WriteString(")")
	}

	b.WriteString(")")

	if _, err := tx.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}

func dropTable(ctx context.Context, tx pgx.Tx, sch *arrow.Schema) error {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("drop table ")
	b.WriteString(tableName)

	if _, err := tx.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}

func alterTable(ctx context.Context, tx pgx.Tx, sch *arrow.Schema, adds, dels []arrow.Field) error {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	css := []string{}
	for _, field := range adds {
		css = append(css, addColumnStatement(&field, true))
	}
	for _, field := range dels {
		css = append(css, dropColumnStatement(&field))
	}

	var b strings.Builder
	b.WriteString("alter table ")
	b.WriteString(tableName)
	b.WriteString(" ")
	b.WriteString(strings.Join(css, ", "))

	if _, err := tx.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}

func dropColumnStatement(f *arrow.Field) string {
	return fmt.Sprintf("drop column %s", sanitize(f.Name))
}

func addColumnStatement(f *arrow.Field, prefix bool) string {
	name := sanitize(f.Name)

	unique := ""
	if metadata.IsUnique(f) {
		unique = "unique"
	}

	null := ""
	if !f.Nullable {
		null = "not null"
	}

	cs := fmt.Sprintf("%s %s %s %s", name, pgarrow.From(f.Type), unique, null)
	if prefix {
		cs = "add column " + cs
	}

	return cs
}

func renewTable(ctx context.Context, tx pgx.Tx, sch *arrow.Schema) error {
	if err := dropTable(ctx, tx, sch); err != nil {
		return err
	}
	return createTableIfNotExists(ctx, tx, sch)
}
