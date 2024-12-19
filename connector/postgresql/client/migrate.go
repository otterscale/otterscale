package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

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

func isKeyChanged(a, b *arrow.Field) bool {
	return a.Nullable != b.Nullable ||
		metadata.IsPrimaryKey(a) != metadata.IsPrimaryKey(b) ||
		metadata.IsUnique(a) != metadata.IsUnique(b)
}

func filterFields(as, bs []arrow.Field) ([]arrow.Field, bool) {
	ret := []arrow.Field{}
	for _, a := range as { // 1, 2, 3, 4, 5
		ok := false
		for _, b := range bs { // 1, 2, 4, 5
			if a.Name != b.Name {
				continue
			}
			if isKeyChanged(&a, &b) {
				return nil, false
			}
			if a.Equal(b) {
				ok = true
				break
			}
		}
		if !ok {
			ret = append(ret, a)
		}
	}
	return ret, true
}

// Only allow adding and deleting fields, but not modifying field names or data types.
func compareSchemata(cur, new *arrow.Schema) (add, del []arrow.Field, ok bool) {
	add, ok = filterFields(cur.Fields(), new.Fields())
	if !ok {
		return
	}
	del, ok = filterFields(new.Fields(), cur.Fields())
	if !ok {
		return
	}
	return
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
	add, del, ok := compareSchemata(cur, new)
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
