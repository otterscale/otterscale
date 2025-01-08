package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (h *Helper) Truncate(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, truncateTableStatement(h.tableName)); err != nil {
		return err
	}
	return nil
}

func truncateTableStatement(tableName string) string {
	var b strings.Builder
	b.WriteString("truncate table ")
	b.WriteString(sanitize(tableName))
	return b.String()
}
