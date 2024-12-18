package client

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/metadata"
)

func (c *Client) toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}

	metadata.SetTableName(m, tableName)
	mtd := arrow.MetadataFrom(m)

	return &mtd
}

func (c *Client) toSchemaFields(header arrow.Record) []arrow.Field {
	var flds []arrow.Field

	if c.opts.infering {
		sch := header.Schema()
		flds := sch.Fields()
		return flds
	} else {
		for _, f := range header.Schema().Fields() {
			flds = append(flds, arrow.Field{Name: f.Name, Type: arrow.BinaryTypes.String})
		}

		return flds
	}
}
