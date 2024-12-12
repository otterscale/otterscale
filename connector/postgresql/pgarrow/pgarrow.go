package pgarrow

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/extensions"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/openhdc/openhdc/internal/codec"
)

var _ codec.Codec = (*Codec)(nil)

type Codec struct {
	codec.DefaultCodec
}

func NewCodec() codec.Codec {
	return &Codec{}
}

func (c *Codec) newArray(dt arrow.DataType, val any) (arrow.Array, error) {
	b := array.NewBuilder(memory.DefaultAllocator, dt)
	defer b.Release()
	if val != nil {
		if err := c.Append(b, val); err != nil {
			return nil, err
		}
	}
	return b.NewArray(), nil
}

//nolint:gocyclo,funlen
func (c *Codec) Encode(typ, val any) (arrow.Array, error) {
	switch typ.(type) {
	case pgtype.BoolCodec:
		return c.newArray(arrow.FixedWidthTypes.Boolean, val)

	case pgtype.ByteaCodec:
		return c.newArray(arrow.BinaryTypes.Binary, val)

	case pgtype.DateCodec:
		return c.newArray(arrow.FixedWidthTypes.Date32, val)

	case pgtype.Float4Codec:
		return c.newArray(arrow.PrimitiveTypes.Float32, val)

	case pgtype.Float8Codec:
		return c.newArray(arrow.PrimitiveTypes.Float64, val)

	case pgtype.Int2Codec:
		return c.newArray(arrow.PrimitiveTypes.Int16, val)

	case pgtype.Int4Codec:
		return c.newArray(arrow.PrimitiveTypes.Int32, val)

	case pgtype.Int8Codec:
		return c.newArray(arrow.PrimitiveTypes.Int64, val)

	case *pgtype.JSONCodec, *pgtype.JSONBCodec:
		typ, _ := extensions.NewJSONType(arrow.BinaryTypes.String)
		return c.newArray(typ, val)

	case pgtype.TextCodec, *pgtype.TextFormatOnlyCodec:
		return c.newArray(arrow.BinaryTypes.String, val)

	case pgtype.TimeCodec, *pgtype.TimestampCodec, *pgtype.TimestamptzCodec:
		return c.newArray(arrow.FixedWidthTypes.Timestamp_us, val)

	case pgtype.Uint32Codec:
		return c.newArray(arrow.PrimitiveTypes.Uint32, val)

	// FIXME: github.com/jackc/pgx/v5 v5.7.1.next
	/*
		case pgtype.Uint64Codec:
			return c.newArray(arrow.PrimitiveTypes.Uint64, val)
	*/

	case pgtype.UUIDCodec:
		return c.newArray(extensions.NewUUIDType(), val)

	// TODO: complete
	case *pgtype.ArrayCodec:
	case pgtype.BitsCodec:
	case pgtype.BoxCodec:
	case pgtype.CircleCodec:
	case pgtype.EnumCodec:
	case pgtype.HstoreCodec:
	case pgtype.InetCodec:
	case pgtype.IntervalCodec:
	case pgtype.LineCodec:
	case pgtype.LsegCodec:
	case pgtype.LtreeCodec:
	case pgtype.MacaddrCodec:
	case *pgtype.MultirangeCodec:
	case pgtype.NumericCodec:
	case pgtype.PathCodec:
	case pgtype.PointCodec:
	case pgtype.PolygonCodec:
	case pgtype.QCharCodec:
	case *pgtype.RangeCodec:
	case pgtype.RecordCodec:
	case pgtype.TIDCodec:
	case *pgtype.XMLCodec:
	}

	return nil, fmt.Errorf("type %T: %w", typ, codec.ErrNotSupported)
}
