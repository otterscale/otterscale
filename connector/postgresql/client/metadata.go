package client

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/metadata"
)

func toFieldMetadata(att *pgAttribute) arrow.Metadata {
	m := map[string]string{}
	if att.ConType != nil && att.ConName != nil {
		switch *att.ConType {
		case 'p':
			metadata.SetPrimaryKey(m)
		case 'u':
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
