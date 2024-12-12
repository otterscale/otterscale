package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func createTableIfNotExists(ctx context.Context, pool *pgxpool.Pool, sch *arrow.Schema) error {
	tableName, ok := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	if !ok {
		return errors.New("table name not found")
	}

	tableName += "_tmp"

	columns := []string{}
	pks := []string{}
	for _, field := range sch.Fields() {
		name := sanitize(field.Name)

		dataType, ok := field.Metadata.GetValue(metadata.KeyFieldDataType)
		if !ok {
			return fmt.Errorf("data type not found: %s", name)
		}

		unique := ""
		if _, ok := field.Metadata.GetValue(metadata.KeyFieldIsUnique); ok {
			unique = "unique"
		}

		null := ""
		if !field.Nullable {
			null = "not null"
		}

		column := fmt.Sprintf("%s %s %s %s", name, dataType, unique, null)
		columns = append(columns, column)

		if _, ok := field.Metadata.GetValue(metadata.KeyFieldIsPrimaryKey); ok {
			pks = append(pks, name)
		}
	}

	var b strings.Builder
	b.WriteString("create table if not exists ")
	b.WriteString(tableName)
	b.WriteString(" (")

	b.WriteString(strings.Join(columns, ", "))

	if len(pks) > 0 {
		b.WriteString(", constraint ")
		b.WriteString(tableName)
		b.WriteString("_pkey")
		b.WriteString(" primary key (")
		b.WriteString(strings.Join(pks, ", "))
		b.WriteString(")")
	}

	b.WriteString(")")

	if _, err := pool.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}

func dropTable(ctx context.Context, pool *pgxpool.Pool, sch *arrow.Schema) error {
	tableName, ok := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	if !ok {
		return errors.New("table name not found")
	}

	var b strings.Builder
	b.WriteString("drop table ")
	b.WriteString(tableName)

	if _, err := pool.Exec(ctx, b.String()); err != nil {
		return err
	}
	return nil
}
