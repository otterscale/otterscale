package pg

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
)

func (h *Helper) Delete(ctx context.Context, tx pgx.Tx, syncedAt time.Time) error {
	if _, err := tx.Exec(ctx, deleteStatement(h.tableName), syncedAt); err != nil {
		return err
	}
	return nil
}

func deleteStatement(tableName string) string {
	var b strings.Builder
	b.WriteString("delete from ")
	b.WriteString(sanitize(tableName))
	b.WriteString(" where ")
	b.WriteString(openhdc.BuiltinFieldSyncedAt())
	b.WriteString(" < $1")
	return b.String()
}
