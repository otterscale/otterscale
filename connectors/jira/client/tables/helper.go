package tables

import (
	"encoding/json"
	"log"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/openhdc/openhdc/metadata"
)

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}
	metadata.SetTableName(m, tableName)
	md := arrow.MetadataFrom(m)
	return &md
}

func builderAppendJson(fbuilder *array.StringBuilder, v any) {
	if v != nil {
		// convert Fields into JSON string
		js, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		fbuilder.Append(string(js))
	} else {
		fbuilder.AppendNull()
	}
}

// type IssueCodec struct{}

// func NewIssueCodec() openhdc.Codec {

// }

// func (c *IssueCodec) Encode(array.Builder, any) error {

// }

// func (c *IssueCodec) Decode(arrow.Array, int) (any, error) {
// 	return nil, openhdc.ErrNotSupported
// }

// var _ openhdc.Codec = (*IssueCodec)(nil)

// func (c *Client) xx(b array.RecordBuilder, sch arrow.Schema, issue *jira.Issue) error {
// 	for idx := range sch.Fields() {
// 		if idx == 0 {
// 			if err := c.Encode(b.Field(idx), issue.Expand); err != nil {
// 				slog.Error("invalid append", "type of field", reflect.TypeOf(b.Field(idx)), "type of value", reflect.TypeOf(val))
// 				return err
// 			}
// 		}

// 	}

// }
