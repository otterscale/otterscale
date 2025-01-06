package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (h *Helper) DropTable(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, dropTableStatement(h.tableName)); err != nil {
		return err
	}
	return nil
}

func dropTableStatement(tableName string) string {
	var b strings.Builder
	b.WriteString("drop table ")
	b.WriteString(sanitize(tableName))
	return b.String()
}
