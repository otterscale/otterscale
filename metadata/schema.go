package metadata

import (
	"errors"

	"github.com/apache/arrow-go/v18/arrow"
)

func SetTableName(m map[string]string, v string) {
	m[keySchemaTableName] = v
}

func GetTableName(s *arrow.Schema) (string, error) {
	v, ok := s.Metadata().GetValue(keySchemaTableName)
	if !ok {
		return "", errors.New("table name not found")
	}
	return v, nil
}
