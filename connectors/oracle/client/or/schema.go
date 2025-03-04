package or

import (
	"context"
	"database/sql/driver"

	go_ora "github.com/sijms/go-ora/v2"
)

type Class struct {
	RelName string `db:"relname"`
}

func Classes(ctx context.Context, pool *go_ora.Connection) ([]*Class, error) {
	sql := `
	SELECT TABLE_NAME FROM USER_TABLES
	`
	stmt := go_ora.NewStmt(sql, pool)
	defer stmt.Close()

	rows, err := stmt.Query_(nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*Class
	for rows.Next_() {
		var tableName string

		if err := rows.Scan(&tableName); err != nil {
			break
		}

		class := &Class{RelName: tableName}

		classes = append(classes, class)
	}

	return classes, nil
}

type Attribute struct {
	AttNum     string  `db:"attnum"`
	AttName    string  `db:"attname"`
	AttTypeID  uint32  `db:"atttypid"`
	AttNotNull bool    `db:"attnotnull"`
	ConTypes   *string `db:"contypes"`
}

func Attributes(ctx context.Context, pool *go_ora.Connection, tableName string) ([]*Attribute, error) {
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
	stmt := go_ora.NewStmt(sql, pool)
	defer stmt.Close()

	rows, err := stmt.Query_([]driver.NamedValue{{Value: tableName}})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oids, err := getTypeOIDs(pool, tableName)
	if err != nil {
		return nil, err
	}

	var attributes []*Attribute
	for rows.Next_() {
		var (
			columnID        string
			columnName      string
			isNullable      bool
			constraintTypes *string
		)

		if err := rows.Scan(&columnID, &columnName, &isNullable, &constraintTypes); err != nil {
			break
		}

		attribute := Attribute{
			AttNum:     columnID,
			AttName:    columnName,
			AttTypeID:  oids[columnID],
			AttNotNull: !isNullable,
			ConTypes:   constraintTypes,
		}

		attributes = append(attributes, &attribute)
	}

	return attributes, nil

}
