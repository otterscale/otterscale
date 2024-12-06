package codec

import (
	"fmt"
	"time"
)

func Format(val any) string {
	switch t := val.(type) {
	case time.Time:
		return t.Format(time.RFC3339Nano)
	default:
		return fmt.Sprintf("%v", val)
	}
}
