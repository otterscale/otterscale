package pg

import (
	"context"
	"log/slog"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc/metadata"
)

func (h *Helper) CreateTable(ctx context.Context, tx pgx.Tx) error {
	slog.Info("[migrate] create table", "table", h.TableName())
	if _, err := tx.Exec(ctx, createTableStatement(h.tableName, h.sch)); err != nil {
		return err
	}
	return nil
}

func createTableStatement(tableName string, sch *arrow.Schema) string {
	ps := []string{}
	ss := []string{}
	for _, f := range sch.Fields() {
		if metadata.IsPrimaryKey(&f) {
			ps = append(ps, sanitize(f.Name))
		}
		ss = append(ss, addColumnStatement(&f, false))
	}

	var b strings.Builder
	b.WriteString("create table if not exists ")
	b.WriteString(sanitize(tableName))
	b.WriteString(" (")
	b.WriteString(strings.Join(ss, ", "))
	if len(ps) > 0 {
		b.WriteString(", constraint ")
		b.WriteString(sanitize(tableName))
		b.WriteString("_pkey")
		b.WriteString(" primary key (")
		b.WriteString(strings.Join(ps, ", "))
		b.WriteString(")")
	}
	b.WriteString(")")
	return b.String()
}
