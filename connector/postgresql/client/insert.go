package client

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc/metadata"
)

func (c *Client) insert(ctx context.Context, rec arrow.Record) error {
	sch := rec.Schema()

	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return err
	}

	columnNames := []string{}
	for _, f := range sch.Fields() {
		columnNames = append(columnNames, f.Name)
	}

	rows := [][]any{}
	for idx := range rec.NumRows() {
		row := []any{}
		for _, col := range rec.Columns() {
			v, err := c.Decode(col, int(idx))
			if err != nil {
				return err
			}
			row = append(row, v)
		}
		rows = append(rows, row)
	}

	_, err = c.pool.CopyFrom(ctx, pgx.Identifier{tableName}, columnNames, pgx.CopyFromRows(rows))
	return err
}
