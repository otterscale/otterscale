package codec

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func format(val any) string {
	switch t := val.(type) {
	case time.Time:
		return t.Format(time.RFC3339Nano)
	case [16]uint8:
		return uuid.UUID(t).String()
	}
	return fmt.Sprintf("%v", val)
}
