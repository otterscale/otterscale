package openhdc

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/metadata"
)

// Only allow adding and deleting fields, but not modifying field names or data types.
func CompareSchemata(cur, new *arrow.Schema) (add, del []arrow.Field, ok bool) {
	add, ok = filterFields(cur.Fields(), new.Fields())
	if !ok {
		return
	}
	del, ok = filterFields(new.Fields(), cur.Fields())
	if !ok {
		return
	}
	return
}

func isFieldKeyChanged(cur, new *arrow.Field) bool {
	return metadata.IsPrimaryKey(cur) != metadata.IsPrimaryKey(new) ||
		metadata.IsUnique(cur) != metadata.IsUnique(new)
}

func filterFields(as, bs []arrow.Field) ([]arrow.Field, bool) {
	diff := []arrow.Field{}
	for _, a := range as {
		if isBuiltinField(&a) {
			continue
		}
		ok := false
		for _, b := range bs {
			if a.Name != b.Name {
				continue
			}
			if isFieldKeyChanged(&a, &b) {
				return nil, false
			}
			if a.Equal(b) {
				ok = true
				break
			}
		}
		if !ok {
			diff = append(diff, a)
		}
	}
	return diff, true
}
