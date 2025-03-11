package or

import (
	"strings"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/metadata"
)

func sanitize(identifier ...string) string {
	parts := make([]string, len(identifier))

	for i := range identifier {
		s := strings.ReplaceAll(identifier[i], string([]byte{0}), "")
		parts[i] = `"` + s + `"`
	}

	return strings.Join(parts, ".")
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
		sch: sch,
		cdc: cdc,

		tableName: tableName,
	}, nil
}

func (h *Helper) TableName() string {
	return h.tableName
}

func (h *Helper) Schema() *arrow.Schema {
	return h.sch
}
