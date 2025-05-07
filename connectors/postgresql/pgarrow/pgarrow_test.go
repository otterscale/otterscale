package pgarrow

import (
	"reflect"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/extensions"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestToForOID(t *testing.T) {
	type args struct {
		oid uint32
	}
	tests := []struct {
		name string
		args args
		want arrow.DataType
	}{
		{
			name: "Known OID",
			args: args{oid: pgtype.TextOID},
			want: arrow.BinaryTypes.String,
		},
		{
			name: "Unknown OID",
			args: args{oid: 99999},
			want: arrow.BinaryTypes.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToForOID(tt.args.oid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToForOID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTo(t *testing.T) {
	type args struct {
		typ string
	}
	tests := []struct {
		name string
		args args
		want arrow.DataType
	}{
		{name: "bytea", args: args{typ: "bytea"}, want: arrow.BinaryTypes.Binary},
		{name: "text", args: args{typ: "text"}, want: arrow.BinaryTypes.String},
		{name: "bool", args: args{typ: "bool"}, want: arrow.FixedWidthTypes.Boolean},
		{name: "date", args: args{typ: "date"}, want: arrow.FixedWidthTypes.Date32},
		{name: "time", args: args{typ: "time"}, want: arrow.FixedWidthTypes.Time64ns},
		{name: "timestamp", args: args{typ: "timestamp"}, want: &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: ""}},
		{name: "timestamptz", args: args{typ: "timestamptz"}, want: arrow.FixedWidthTypes.Timestamp_us},
		{name: "float4", args: args{typ: "float4"}, want: arrow.PrimitiveTypes.Float32},
		{name: "float8", args: args{typ: "float8"}, want: arrow.PrimitiveTypes.Float64},
		{name: "int2", args: args{typ: "int2"}, want: arrow.PrimitiveTypes.Int16},
		{name: "int4", args: args{typ: "int4"}, want: arrow.PrimitiveTypes.Int32},
		{name: "int8", args: args{typ: "int8"}, want: arrow.PrimitiveTypes.Int64},
		{name: "uuid", args: args{typ: "uuid"}, want: extensions.NewUUIDType()},
		{name: "unknown type", args: args{typ: "unknown_type"}, want: arrow.BinaryTypes.String},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := To(tt.args.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	type args struct {
		typ arrow.DataType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "BinaryType", args: args{typ: arrow.BinaryTypes.Binary}, want: "bytea"},
		{name: "StringType", args: args{typ: arrow.BinaryTypes.String}, want: "text"},
		{name: "BooleanType", args: args{typ: arrow.FixedWidthTypes.Boolean}, want: "bool"},
		{name: "Date32Type", args: args{typ: arrow.FixedWidthTypes.Date32}, want: "date"},
		{name: "Time64Type", args: args{typ: arrow.FixedWidthTypes.Time64ns}, want: "time"},
		{name: "TimestampType without timezone", args: args{typ: &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: ""}}, want: "timestamp"},
		{name: "TimestampType with timezone", args: args{typ: &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: "UTC"}}, want: "timestamptz"},
		{name: "Float32Type", args: args{typ: arrow.PrimitiveTypes.Float32}, want: "float4"},
		{name: "Float64Type", args: args{typ: arrow.PrimitiveTypes.Float64}, want: "float8"},
		{name: "Int16Type", args: args{typ: arrow.PrimitiveTypes.Int16}, want: "int2"},
		{name: "Int32Type", args: args{typ: arrow.PrimitiveTypes.Int32}, want: "int4"},
		{name: "Int64Type", args: args{typ: arrow.PrimitiveTypes.Int64}, want: "int8"},
		{name: "Uint32Type", args: args{typ: arrow.PrimitiveTypes.Uint32}, want: "xid"},
		{name: "Uint64Type", args: args{typ: arrow.PrimitiveTypes.Uint64}, want: "xid8"},
		{name: "UUIDType", args: args{typ: extensions.NewUUIDType()}, want: "uuid"},
		{name: "UnknownType", args: args{typ: arrow.BinaryTypes.LargeBinary}, want: "text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From(tt.args.typ); got != tt.want {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}
}
