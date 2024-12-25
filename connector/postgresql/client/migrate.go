package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/metadata"
)

func getTable(schs []*arrow.Schema, new string) (*arrow.Schema, bool) {
	for _, sch := range schs {
		current, _ := metadata.GetTableName(sch)
		if current == new {
			return sch, true
		}
	}
	return nil, false
}

func migrate(ctx context.Context, tx pgx.Tx, schs []*arrow.Schema, new *arrow.Schema) error {
	tableName, err := metadata.GetTableName(new)
	if err != nil {
		return err
	}
	cur, ok := getTable(schs, tableName)
	if !ok {
		fmt.Printf("[migrate] create table %s\n", tableName)
		return createTableIfNotExists(ctx, tx, new)
	}
	add, del, ok := openhdc.CompareSchemata(cur, new)
	if !ok {
		fmt.Printf("[migrate] renew table %s\n", tableName)
		return renewTable(ctx, tx, new)
	}
	if len(add) > 0 || len(del) > 0 {
		fmt.Printf("[migrate] alter table %s: add %v del %v\n", tableName, add, del)
		if err := alterTable(ctx, tx, new, add, del); err != nil {
			slog.Error(err.Error())
			fmt.Printf("[migrate] renew table %s\n", tableName)
			return renewTable(ctx, tx, new)
		}
	}
	fmt.Printf("[migrate] skip table %s\n", tableName)
	return nil
}
