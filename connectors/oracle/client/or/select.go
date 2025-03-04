package or

import (
	"context"
	"database/sql/driver"
	"fmt"
	"slices"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	go_ora "github.com/sijms/go-ora/v2"

	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

func (h *Helper) Select(ctx context.Context, pool *go_ora.Connection, mode property.SyncMode, curs []*workload.Sync_Option_Cursor) (driver.Rows, error) {
	stmt := go_ora.NewStmt(selectStatement(h.tableName, h.sch, mode, curs), pool)
	defer stmt.Close()

	rows, err := stmt.Query(nil)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func selectStatement(tableName string, sch *arrow.Schema, mode property.SyncMode, curs []*workload.Sync_Option_Cursor) string {
	cs := []string{}
	for _, f := range sch.Fields() {
		cs = append(cs, sanitize(f.Name))
	}

	var b strings.Builder
	b.WriteString("select ")
	b.WriteString(strings.Join(cs, ", "))
	b.WriteString(" from ")
	b.WriteString(sanitize(tableName))
	b.WriteString(cursorsToWhere(mode, curs))
	return b.String()
}

func cursorsToWhere(mode property.SyncMode, curs []*workload.Sync_Option_Cursor) string {
	supports := []property.SyncMode{
		property.SyncMode_incremental_append,
		property.SyncMode_incremental_append_dedupe,
	}
	if len(curs) == 0 || !slices.Contains(supports, mode) {
		return ""
	}
	ws := []string{}
	for _, cur := range curs {
		ws = append(ws, fmt.Sprintf("%s > '%s'", sanitize(cur.GetField()), cur.GetValue()))
	}
	b := strings.Builder{}
	b.WriteString(" where ")
	b.WriteString(strings.Join(ws, " and "))
	return b.String()
}
