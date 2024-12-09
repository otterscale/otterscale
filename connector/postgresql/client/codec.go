package client

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/openhdc/openhdc/internal/codec"
)

var _ codec.Codec = (*postgresqlCodec)(nil)

type postgresqlCodec struct {
	codec.DefaultCodec
}

func NewCodec() codec.Codec {
	return &postgresqlCodec{}
}

func (c *postgresqlCodec) newArray(dt arrow.DataType, val any) (arrow.Array, error) {
	b := array.NewBuilder(memory.DefaultAllocator, dt)
	defer b.Release()
	if val != nil {
		if err := c.Append(b, val); err != nil {
			return nil, err
		}
	}
	return b.NewArray(), nil
}

//nolint:gocyclo
func (c *postgresqlCodec) Encode(typ, val any) (arrow.Array, error) {
	switch typ.(type) {
	case *pgtype.ArrayCodec:
	case pgtype.BitsCodec:
	case pgtype.BoolCodec:
		return c.newArray(arrow.FixedWidthTypes.Boolean, val)

	case pgtype.BoxCodec:
	case pgtype.ByteaCodec:
	case pgtype.CircleCodec:
	case *pgtype.CompositeCodec:
	case pgtype.DateCodec:
	case *pgtype.DateCodec:
	case *pgtype.EnumCodec:
	case pgtype.Float4Codec:
	case pgtype.Float8Codec:
	case pgtype.HstoreCodec:
	case pgtype.InetCodec:
	case pgtype.Int2Codec:
	case pgtype.Int4Codec:
	case pgtype.Int8Codec:
		return c.newArray(arrow.PrimitiveTypes.Int64, val)

	case pgtype.IntervalCodec:
	case *pgtype.JSONCodec:
	case *pgtype.JSONBCodec:
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
	case pgtype.TextCodec:
		return c.newArray(arrow.BinaryTypes.String, val)

	case *pgtype.TextFormatOnlyCodec:
	case pgtype.TIDCodec:
	case pgtype.TimeCodec:
	case *pgtype.TimestampCodec:
	case *pgtype.TimestamptzCodec:
		return c.newArray(arrow.FixedWidthTypes.Timestamp_us, val)

	case pgtype.Uint32Codec:
	// FIXME: github.com/jackc/pgx/v5 v5.7.1.next
	// case pgtype.Uint64Codec:
	case pgtype.UUIDCodec:
	case *pgtype.XMLCodec:
	}

	return nil, fmt.Errorf("type %T: %w", typ, codec.ErrNotSupported)
}
