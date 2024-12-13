package client

import (
	"context"
	"errors"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/metadata"
)

func getTable(tabs []arrow.Table, new string) (*arrow.Schema, bool) {
	for _, tab := range tabs {
		sch := tab.Schema()
		current, _ := sch.Metadata().GetValue(metadata.KeySchemaTableName)
		if current == new {
			return sch, true
		}
	}
	return nil, false
}

func isKeyChanged(a, b *arrow.Field) bool {
	apk, _ := a.Metadata.GetValue(metadata.KeyFieldIsPrimaryKey)
	auq, _ := a.Metadata.GetValue(metadata.KeyFieldIsUnique)
	bpk, _ := b.Metadata.GetValue(metadata.KeyFieldIsPrimaryKey)
	buq, _ := b.Metadata.GetValue(metadata.KeyFieldIsUnique)
	return a.Nullable != b.Nullable || apk != bpk || auq != buq
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

func (c *Client) migrate(ctx context.Context, tabs []arrow.Table, sch *arrow.Schema) error {
	tableName, ok := sch.Metadata().GetValue(metadata.KeySchemaTableName)
	if !ok {
		return errors.New("table name not found")
	}
	cur, ok := getTable(tabs, tableName)
	if !ok {
		return createTableIfNotExists(ctx, c.pool, sch)
	}
	add, del, ok := compareSchemata(cur, sch)
	if !ok {
		return renewTable(ctx, c.pool, sch)
	}
	if len(add) > 0 || len(del) > 0 {
		if err := alterTable(ctx, c.pool, sch, add, del); err != nil {
			slog.Error(err.Error())
			return renewTable(ctx, c.pool, sch)
		}
	}
	return nil
}
