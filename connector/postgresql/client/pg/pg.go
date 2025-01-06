package pg

import (
	"context"
	"log/slog"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jackc/pgx/v5"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/metadata"
)

func sanitize(str string) string {
	return pgx.Identifier{str}.Sanitize()
}

type Helper struct {
	sch *arrow.Schema
	cdc openhdc.Codec

	tableName string
}

func NewHelper(sch *arrow.Schema, cdc openhdc.Codec) (*Helper, error) {
	tableName, err := metadata.GetTableName(sch)
	if err != nil {
		return nil, err
	}
	return &Helper{
		sch:       sch,
		cdc:       cdc,
		tableName: tableName,
	}, nil
}

func (h *Helper) TableName() string {
	return h.tableName
}

func (h *Helper) Schema() *arrow.Schema {
	return h.sch
}

func (h *Helper) RenewTable(ctx context.Context, tx pgx.Tx) error {
	slog.Info("[migrate] renew table", "table", h.TableName())
	if err := h.DropTable(ctx, tx); err != nil {
		return err
	}
	return h.CreateTable(ctx, tx)
}
