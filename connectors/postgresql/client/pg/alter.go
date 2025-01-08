package pg

import (
	"context"
	"log/slog"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc/connectors/postgresql/pgarrow"
	"github.com/openhdc/openhdc/metadata"
)

func (h *Helper) AlterTable(ctx context.Context, tx pgx.Tx, adds, dels []arrow.Field) error {
	slog.Info("[migrate] alter table", "table", h.TableName(), "add", adds, "del", dels)
	if _, err := tx.Exec(ctx, alterTableStatement(h.tableName, adds, dels)); err != nil {
		return err
	}
	return nil
}

func alterTableStatement(tableName string, adds, dels []arrow.Field) string {
	ss := []string{}
	for _, f := range adds {
		ss = append(ss, addColumnStatement(&f, true))
	}
	for _, f := range dels {
		ss = append(ss, dropColumnStatement(&f))
	}

	var b strings.Builder
	b.WriteString("alter table ")
	b.WriteString(sanitize(tableName))
	b.WriteString(" ")
	b.WriteString(strings.Join(ss, ", "))
	return b.String()
}

func dropColumnStatement(f *arrow.Field) string {
	var b strings.Builder
	b.WriteString("drop column ")
	b.WriteString(sanitize(f.Name))
	return b.String()
}

func addColumnStatement(f *arrow.Field, prefix bool) string {
	unique := ""
	if metadata.IsUnique(f) {
		unique = "unique"
	}

	null := ""
	if !f.Nullable {
		null = "not null"
	}

	var b strings.Builder
	if prefix {
		b.WriteString("add column ")
	}
	b.WriteString(sanitize(f.Name))
	b.WriteString(" ")
	b.WriteString(pgarrow.From(f.Type))
	b.WriteString(" ")
	b.WriteString(unique)
	b.WriteString(" ")
	b.WriteString(null)
	return b.String()
}
