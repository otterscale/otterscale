package client

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/internal/metadata"
)

func toFieldMetadata(att *pgAttribute, dataType string) arrow.Metadata {
	m := map[string]string{
		metadata.KeyFieldDataType: dataType,
	}
	if att.ConType != nil && att.ConName != nil {
		switch *att.ConType {
		case 'p':
			m[metadata.KeyFieldIsPrimaryKey] = "true"
		case 'u':
			m[metadata.KeyFieldIsUnique] = "true"
		}
	}
	return arrow.MetadataFrom(m)
}

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{
		metadata.KeySchemaTableName: tableName,
	}
	md := arrow.MetadataFrom(m)
	return &md
}
