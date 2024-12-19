package client

import (
	"strings"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/metadata"
)

func toFieldMetadata(att *pgAttribute) arrow.Metadata {
	m := map[string]string{}
	if att.ConTypes != nil {
		if strings.Contains(*att.ConTypes, "p") {
			metadata.SetPrimaryKey(m)
		}
		if strings.Contains(*att.ConTypes, "u") {
			metadata.SetUnique(m)
		}
	}
	return arrow.MetadataFrom(m)
}

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}
	metadata.SetTableName(m, tableName)
	md := arrow.MetadataFrom(m)
	return &md
}
