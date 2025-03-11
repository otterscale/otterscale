package orarrow

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/extensions"

	go_ora "github.com/sijms/go-ora/v2"
)

func ToForOID(oid uint32) arrow.DataType {
	typ_name := go_ora.TNSType(oid).String()
	return To(typ_name)
}

func To(typ string) arrow.DataType {
	switch typ {
	case "NCHAR":
		return arrow.BinaryTypes.String
	case "ROWID", "UROWID":
		return arrow.BinaryTypes.String
	case "LONG":
		return arrow.BinaryTypes.LargeString
	case "OCIClobLocator":
		return arrow.BinaryTypes.LargeString
	case "NUMBER":
		return arrow.PrimitiveTypes.Float64
	case "IBFloat":
		return arrow.PrimitiveTypes.Float32
	case "IBDouble":
		return arrow.PrimitiveTypes.Float64
	case "RAW":
		return arrow.BinaryTypes.Binary
	case "LongRaw":
		return arrow.BinaryTypes.LargeBinary
	case "OCIBlobLocator", "OCIFileLocator":
		return arrow.BinaryTypes.LargeBinary
	case "DATE":
		return arrow.FixedWidthTypes.Date32
	case "TimeStampDTY":
		return &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: ""}
	case "TimeStampTZ_DTY", "TimeStampLTZ_DTY":
		return arrow.FixedWidthTypes.Timestamp_us
	case "IntervalYM_DTY":
		return &arrow.MonthIntervalType{}
	case "IntervalDS_DTY":
		return &arrow.DurationType{Unit: arrow.TimeUnit(arrow.Second)}
	case "JSON":
		return &extensions.JSONType{}
	}
	return arrow.BinaryTypes.String
}

func From(typ arrow.DataType) string {
	return "LONG"
}
