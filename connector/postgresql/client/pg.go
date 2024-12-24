package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc"
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

func truncate(ctx context.Context, tx pgx.Tx, sch *arrow.Schema) error {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("truncate table ")
	b.WriteString(tableName)

	if _, err := tx.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}

func delete(ctx context.Context, tx pgx.Tx, rec arrow.Record, syncedAt time.Time) error {
	tableName, err := metadata.GetTableName(rec.Schema())
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("delete from ")
	b.WriteString(tableName)
	b.WriteString(" where _openhdc_synced_at < $1")

	if _, err := tx.Exec(ctx, b.String(), syncedAt); err != nil {
		return err
	}
	return nil
}

func insertStatement(rec arrow.Record, table string) string {
	columns := []string{}
	values := []string{}
	for i, f := range rec.Schema().Fields() {
		columns = append(columns, sanitize(f.Name))
		values = append(values, fmt.Sprintf("$%d", i+1))
	}

	var b strings.Builder
	b.WriteString("insert into ")
	b.WriteString(table)
	b.WriteString(" (")
	b.WriteString(strings.Join(columns, ", "))
	b.WriteString(") values (")
	b.WriteString(strings.Join(values, ", "))
	b.WriteString(")")

	return b.String()
}

// auto combine as batch in transaction
func insert(ctx context.Context, tx pgx.Tx, rec arrow.Record, c openhdc.Codec) error {
	table, err := metadata.GetTableName(rec.Schema())
	if err != nil {
		return err
	}
	args := []any{}
	for _, col := range rec.Columns() {
		v, err := c.Decode(col, 0)
		if err != nil {
			return err
		}
		args = append(args, v)
	}
	if _, err := tx.Exec(ctx, insertStatement(rec, table), args...); err != nil {
		return err
	}
	return nil
}

func upsertStatement(rec arrow.Record, table string, update bool) string {
	var b strings.Builder
	b.WriteString(insertStatement(rec, table))
	if !update {
		b.WriteString(" on conflict do nothing")
		return b.String()
	}
	pks := []string{}
	columns := []string{}
	for _, f := range rec.Schema().Fields() {
		if metadata.IsPrimaryKey(&f) {
			pks = append(pks, sanitize(f.Name))
			continue
		}
		columns = append(columns, fmt.Sprintf("%[1]s = excluded.%[1]s", sanitize(f.Name)))
	}
	b.WriteString(" on conflict (")
	b.WriteString(strings.Join(pks, ", "))
	b.WriteString(") do ")
	if len(columns) == 0 {
		b.WriteString("nothing")
		return b.String()
	}
	b.WriteString("update set ")
	b.WriteString(strings.Join(columns, ", "))
	return b.String()
}

func upsert(ctx context.Context, tx pgx.Tx, rec arrow.Record, c openhdc.Codec, update bool) error {
	table, err := metadata.GetTableName(rec.Schema())
	if err != nil {
		return err
	}
	args := []any{}
	for _, col := range rec.Columns() {
		v, err := c.Decode(col, 0)
		if err != nil {
			return err
		}
		args = append(args, v)
	}
	x := upsertStatement(rec, table, update)
	fmt.Println(x)
	if _, err := tx.Exec(ctx, x, args...); err != nil {
		return err
	}
	return nil
}
