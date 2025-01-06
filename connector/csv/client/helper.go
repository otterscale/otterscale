package client

import (
	"errors"
	"slices"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/metadata"
)

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}
	metadata.SetTableName(m, tableName)
	md := arrow.MetadataFrom(m)
	return &md
}

func toStringFields(fs []arrow.Field) []arrow.Field {
	nfs := []arrow.Field{}
	for _, f := range fs {
		nfs = append(nfs, arrow.Field{Name: f.Name, Type: arrow.BinaryTypes.String})
	}
	return nfs
}

func validateFields(fs []arrow.Field) error {
	if slices.ContainsFunc(fs, func(f arrow.Field) bool { return f.Name == "" }) {
		return errors.New("header row has empty value")
	}
	cfs := slices.CompactFunc(fs, func(a, b arrow.Field) bool { return a.Name == b.Name })
	if len(fs) != len(cfs) {
		return errors.New("header row has duplicate value")
	}
	return nil
}
