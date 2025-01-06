package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc/metadata"
)

func (h *Helper) Upsert(ctx context.Context, tx pgx.Tx, rows int64, cols []arrow.Array, update bool) error {
	b := &pgx.Batch{}
	for row := range rows {
		args := []any{}
		for _, col := range cols {
			v, err := h.cdc.Decode(col, int(row))
			if err != nil {
				return err
			}
			args = append(args, v)
		}
		b.Queue(upsertStatement(h.tableName, h.sch, update), args...)
	}
	return tx.SendBatch(ctx, b).Close()
}

func upsertStatement(tableName string, sch *arrow.Schema, update bool) string {
	var b strings.Builder
	b.WriteString(insertStatement(tableName, sch))
	if !update {
		b.WriteString(" on conflict do nothing")
		return b.String()
	}
	ps := []string{}
	cs := []string{}
	for _, f := range sch.Fields() {
		if metadata.IsPrimaryKey(&f) {
			ps = append(ps, sanitize(f.Name))
			continue
		}
		cs = append(cs, fmt.Sprintf("%[1]s = excluded.%[1]s", sanitize(f.Name)))
	}
	b.WriteString(" on conflict (")
	b.WriteString(strings.Join(ps, ", "))
	b.WriteString(") do ")
	if len(cs) == 0 {
		b.WriteString("nothing")
		return b.String()
	}
	b.WriteString("update set ")
	b.WriteString(strings.Join(cs, ", "))
	return b.String()
}
