package openhdc

import (
	"slices"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/metadata"
)

func isBuiltin(a *arrow.Field) bool {
	return slices.Contains([]string{
		BuiltinFieldName,
		BuiltinFieldSyncedAt,
	}, a.Name)
}

func isKeyChanged(a, b *arrow.Field) bool {
	return a.Nullable != b.Nullable ||
		metadata.IsPrimaryKey(a) != metadata.IsPrimaryKey(b) ||
		metadata.IsUnique(a) != metadata.IsUnique(b)
}

func filterFields(as, bs []arrow.Field) ([]arrow.Field, bool) {
	ret := []arrow.Field{}
	for _, a := range as { // 1, 2, 3, 4, 5
		if isBuiltin(&a) {
			continue
		}
		ok := false
		for _, b := range bs { // 1, 2, 4, 5
			if a.Name != b.Name {
				continue
			}
			if isKeyChanged(&a, &b) {
				return nil, false
			}
			if a.Equal(b) {
				ok = true
				break
			}
		}
		if !ok {
			ret = append(ret, a)
		}
	}
	return ret, true
}

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
