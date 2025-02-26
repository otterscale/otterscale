package or

import (
	"context"
	"database/sql"
)

type Class struct {
	RelName string `db:"relname"`
}

func readTables(ctx context.Context, pool *sql.DB) (*sql.Rows, error) {
	sql := `
		SELECT TABLE_NAME FROM USER_TABLES
	`
	rows, err := pool.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func Classes(ctx context.Context, pool *sql.DB) ([]*Class, error) {
	var tables []*Class

	rows, err := readTables(ctx, pool)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table string

		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, &Class{RelName: table})
	}

	return tables, nil
}

type Attribute struct {
	AttNum     string  `db:"attnum"`
	AttName    string  `db:"attname"`
	AttTypeID  uint32  `db:"atttypid"`
	AttNotNull bool    `db:"attnotnull"`
	ConTypes   *string `db:"contypes"`
}

func readColumns(ctx context.Context, pool *sql.DB, tableName string) (*sql.Rows, error) {
	sql := `
		WITH CONS_TYPES AS (
			SELECT
				cc.TABLE_NAME, cc.COLUMN_NAME, LISTAGG(c.CONSTRAINT_TYPE, ',') CONSTRAINT_TYPES
			FROM
				USER_CONSTRAINTS c
				LEFT JOIN
				USER_CONS_COLUMNS cc ON c.TABLE_NAME = cc.TABLE_NAME AND c.CONSTRAINT_NAME = cc.CONSTRAINT_NAME
			GROUP BY
				cc.TABLE_NAME, cc.COLUMN_NAME	
		)
		SELECT
			tc.COLUMN_ID, tc.COLUMN_NAME, TO_BOOLEAN(tc.NULLABLE) NULLABLE, ct.CONSTRAINT_TYPES
		FROM
			USER_TAB_COLUMNS tc
			LEFT JOIN
			CONS_TYPES ct ON tc.TABLE_NAME = ct.TABLE_NAME AND tc.COLUMN_NAME = ct.COLUMN_NAME
		WHERE
			tc.TABLE_NAME = :1
	`
	rows, err := pool.QueryContext(ctx, sql, tableName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func Attributes(ctx context.Context, pool *sql.DB, tableName string) ([]*Attribute, error) {
	var attributes []*Attribute

	oids, err := getTypeOIDs(pool, tableName)
	if err != nil {
		return nil, err
	}

	rows, err := readColumns(ctx, pool, tableName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			column_id        string
			column_name      string
			nullable         bool
			constraint_types *string
		)

		if err := rows.Scan(&column_id, &column_name, &nullable, &constraint_types); err != nil {
			return nil, err
		}

		attribute := Attribute{
			AttNum:     column_id,
			AttName:    column_name,
			AttTypeID:  oids[column_id],
			AttNotNull: !nullable,
			ConTypes:   constraint_types,
		}

		attributes = append(attributes, &attribute)
	}

	return attributes, nil
}
