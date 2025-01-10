package client

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

