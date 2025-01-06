package metadata

import "github.com/apache/arrow-go/v18/arrow"

func SetPrimaryKey(m map[string]string) {
	m[keyFieldIsPrimaryKey] = "true"
}

func IsPrimaryKey(f *arrow.Field) bool {
	_, ok := f.Metadata.GetValue(keyFieldIsPrimaryKey)
	return ok
}

func HasPrimaryKey(s *arrow.Schema) bool {
	for _, f := range s.Fields() {
		if IsPrimaryKey(&f) {
			return true
		}
	}
	return false
}

func SetUnique(m map[string]string) {
	m[keyFieldIsUnique] = "true"
}

func IsUnique(f *arrow.Field) bool {
	_, ok := f.Metadata.GetValue(keyFieldIsUnique)
	return ok
}
