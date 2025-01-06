package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"
)

// auto combine as batch in transaction
func (h *Helper) Insert(ctx context.Context, tx pgx.Tx, rows int64, cols []arrow.Array) error {
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
		b.Queue(insertStatement(h.tableName, h.sch), args...)
	}
	return tx.SendBatch(ctx, b).Close()
}

func insertStatement(tableName string, sch *arrow.Schema) string {
	cs := []string{}
	vs := []string{}
	for i, f := range sch.Fields() {
		cs = append(cs, sanitize(f.Name))
		vs = append(vs, fmt.Sprintf("$%d", i+1))
	}

	var b strings.Builder
	b.WriteString("insert into ")
	b.WriteString(sanitize(tableName))
	b.WriteString(" (")
	b.WriteString(strings.Join(cs, ", "))
	b.WriteString(") values (")
	b.WriteString(strings.Join(vs, ", "))
	b.WriteString(")")
	return b.String()
}
