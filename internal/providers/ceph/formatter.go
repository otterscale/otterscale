package ceph

import (
	"fmt"
	"strconv"
)

func parseQuota(v any) (uint64, error) {
	switch val := v.(type) {
	case string:
		return strconv.ParseUint(val, 10, 64)
	case float64:
		return uint64(val), nil
	case int:
		if val < 0 {
			return 0, fmt.Errorf("quota value is negative: %d", val)
		}
		return uint64(val), nil
	case int64:
		if val < 0 {
			return 0, fmt.Errorf("quota value is negative: %d", val)
		}
		return uint64(val), nil
	case uint64:
		return val, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported quota type %T", v)
	}
}
