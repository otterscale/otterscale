package openhdc

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/extensions"
	"github.com/google/uuid"
)

type Codec interface {
	Encode(b array.Builder, val any) error
	Decode(arr arrow.Array, idx int) (any, error)
}

type DefaultCodec struct{}

var _ Codec = (*DefaultCodec)(nil)

func NewDefaultCodec() Codec {
	return &DefaultCodec{}
}

func (DefaultCodec) format(val any) string {
	switch t := val.(type) {
	case time.Time:
		return t.Format(time.RFC3339Nano)
	case []uint8:
		return base64.StdEncoding.EncodeToString(t)
	case [16]uint8:
		return uuid.UUID(t).String()
	}
	return fmt.Sprintf("%v", val)
}

func (c DefaultCodec) Encode(b array.Builder, val any) error {
	if val == nil {
		b.AppendNull()
		return nil
	}
	return b.AppendValueFromString(c.format(val))
}

func (DefaultCodec) Decode(arr arrow.Array, idx int) (any, error) { //nolint:funlen,gocyclo
	switch a := arr.(type) {
	case *array.Null:
		return a.Value(idx), nil
	case *array.Boolean:
		return a.Value(idx), nil
	case *array.Uint8:
		return a.Value(idx), nil
	case *array.Int8:
		return a.Value(idx), nil
	case *array.Uint16:
		return a.Value(idx), nil
	case *array.Int16:
		return a.Value(idx), nil
	case *array.Uint32:
		return a.Value(idx), nil
	case *array.Int32:
		return a.Value(idx), nil
	case *array.Uint64:
		return a.Value(idx), nil
	case *array.Int64:
		return a.Value(idx), nil
	case *array.Float16:
		return a.Value(idx), nil
	case *array.Float32:
		return a.Value(idx), nil
	case *array.Float64:
		return a.Value(idx), nil
	case *array.String:
		return a.Value(idx), nil
	case *array.Binary:
		return a.Value(idx), nil
	case *array.FixedSizeBinary:
		return a.Value(idx), nil
	case *array.Date32:
		return a.Value(idx).ToTime(), nil
	case *array.Date64:
		return a.Value(idx).ToTime(), nil
	case *array.Timestamp:
		t := arr.DataType().(*arrow.TimestampType)
		return a.Value(idx).ToTime(t.TimeUnit()), nil
	case *array.Time32:
		t := arr.DataType().(*arrow.Time32Type)
		return a.Value(idx).ToTime(t.TimeUnit()), nil
	case *array.Time64:
		t := arr.DataType().(*arrow.Time64Type)
		return a.Value(idx).ToTime(t.TimeUnit()), nil
	case *array.MonthInterval:
		return time.Duration(a.Value(idx)) * time.Hour * 24 * 30, nil //nolint:mnd
	case *array.DayTimeInterval:
		days := time.Duration(a.Value(idx).Days) * time.Hour * 24
		ms := time.Duration(a.Value(idx).Milliseconds) * time.Millisecond
		return days + ms, nil
	case *array.Decimal128:
		t := arr.DataType().(*arrow.Decimal128Type)
		return a.Value(idx).ToFloat64(t.Scale), nil
	case *array.Decimal256:
		t := arr.DataType().(*arrow.Decimal256Type)
		return a.Value(idx).ToFloat64(t.Scale), nil
	case *array.List:
		// FIXME: complete
	case *array.Struct:
		// FIXME: complete
	case *array.SparseUnion:
		// FIXME: complete
	case *array.DenseUnion:
		// FIXME: complete
	case *array.Dictionary:
		// FIXME: complete
	case *array.Map:
		// FIXME: complete
	case *array.FixedSizeList:
		// FIXME: complete
	case *array.Duration:
		t := arr.DataType().(*arrow.DurationType)
		return time.Duration(a.Value(idx)) * t.Unit.Multiplier(), nil
	case *array.LargeString:
		return a.Value(idx), nil
	case *array.LargeBinary:
		return a.Value(idx), nil
	case *array.LargeList:
		// FIXME: complete
	case *array.MonthDayNanoInterval:
		months := time.Duration(a.Value(idx).Months) * time.Hour * 24 * 30
		days := time.Duration(a.Value(idx).Days) * time.Hour * 24
		ns := time.Duration(a.Value(idx).Nanoseconds) * time.Nanosecond
		return months + days + ns, nil
	case *array.RunEndEncoded:
		// FIXME: complete
	case *array.StringView:
		return a.Value(idx), nil
	case *array.BinaryView:
		return a.Value(idx), nil
	case *array.ListView:
		// FIXME: complete
	case *array.LargeListView:
		// FIXME: complete
	case *array.Decimal32:
		t := arr.DataType().(*arrow.Decimal32Type)
		return a.Value(idx).ToFloat64(t.Scale), nil
	case *array.Decimal64:
		t := arr.DataType().(*arrow.Decimal64Type)
		return a.Value(idx).ToFloat64(t.Scale), nil
	case *extensions.Bool8Array:
		return a.Value(idx), nil
	case *extensions.JSONArray:
		return a.Value(idx), nil
	case *extensions.OpaqueArray:
		return a.ValueStr(idx), nil
	case *extensions.UUIDArray:
		return a.Value(idx), nil
	}
	return nil, fmt.Errorf("type %T: %w", arr, ErrNotSupported)
}
