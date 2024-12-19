package pgarrow

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/extensions"
	"github.com/jackc/pgx/v5/pgtype"
)

var typeMap = pgtype.NewMap()

func ToForOID(oid uint32) arrow.DataType {
	if t, ok := typeMap.TypeForOID(oid); ok {
		return To(t.Name)
	}
	return arrow.BinaryTypes.String
}

func To(typ string) arrow.DataType {
	switch typ {
	case "bytea":
		return arrow.BinaryTypes.Binary
	case "aclitem", "bpchar", "jsonpath", "name", "text", "unknown", "varchar", "xml":
		return arrow.BinaryTypes.String
	case "json", "jsonb":
		return arrow.BinaryTypes.String
	case "bool":
		return arrow.FixedWidthTypes.Boolean
	case "date":
		return arrow.FixedWidthTypes.Date32
	case "time":
		return arrow.FixedWidthTypes.Time64ns
	case "timestamp":
		return &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: ""}
	case "timestamptz":
		return arrow.FixedWidthTypes.Timestamp_us
	case "float4":
		return arrow.PrimitiveTypes.Float32
	case "float8":
		return arrow.PrimitiveTypes.Float64
	case "int2":
		return arrow.PrimitiveTypes.Int16
	case "int4":
		return arrow.PrimitiveTypes.Int32
	case "int8":
		return arrow.PrimitiveTypes.Int64
	case "cid", "oid", "xid":
		return arrow.PrimitiveTypes.Uint32
	case "xid8":
		return arrow.PrimitiveTypes.Uint64
	case "uuid":
		return extensions.NewUUIDType()
	}
	return arrow.BinaryTypes.String
}

func From(typ arrow.DataType) string {
	switch v := typ.(type) {
	case *arrow.BinaryType:
		return "bytea"
	case *arrow.StringType:
		return "text"
	case *arrow.BooleanType:
		return "bool"
	case *arrow.Date32Type:
		return "date"
	case *arrow.Time64Type:
		return "time"
	case *arrow.TimestampType:
		if v.TimeZone == "" {
			return "timestamp"
		}
		return "timestamptz"
	case *arrow.Float32Type:
		return "float4"
	case *arrow.Float64Type:
		return "float8"
	case *arrow.Int16Type:
		return "int2"
	case *arrow.Int32Type:
		return "int4"
	case *arrow.Int64Type:
		return "int8"
	case *arrow.Uint32Type:
		return "xid"
	case *arrow.Uint64Type:
		return "xid8"
	case *extensions.UUIDType:
		return "uuid"
	}
	return "text"
}
