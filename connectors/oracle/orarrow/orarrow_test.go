package orarrow

import (
	"reflect"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	// "go get github.com/openhdc/openhdc/connectors/oracle/orarrow"
)

func TestTo(t *testing.T) {
	type args struct {
		typ string
	}
	tests := []struct {
		name string
		args args
		want arrow.DataType
	}{
		{
			name: "Test NCHAR",
			args: args{typ: "NCHAR"},
			want: arrow.BinaryTypes.String,
		},
		{
			name: "Test ROWID",
			args: args{typ: "ROWID"},
			want: arrow.BinaryTypes.String,
		},
		{
			name: "Test LONG",
			args: args{typ: "LONG"},
			want: arrow.BinaryTypes.LargeString,
		},
		{
			name: "Test OCIClobLocator",
			args: args{typ: "OCIClobLocator"},
			want: arrow.BinaryTypes.LargeString,
		},
		{
			name: "Test NUMBER",
			args: args{typ: "NUMBER"},
			want: arrow.PrimitiveTypes.Float64,
		},
		{
			name: "Test IBFloat",
			args: args{typ: "IBFloat"},
			want: arrow.PrimitiveTypes.Float32,
		},
		{
			name: "Test IBDouble",
			args: args{typ: "IBDouble"},
			want: arrow.PrimitiveTypes.Float64,
		},
		{
			name: "Test RAW",
			args: args{typ: "RAW"},
			want: arrow.BinaryTypes.Binary,
		},
		{
			name: "Test LongRaw",
			args: args{typ: "LongRaw"},
			want: arrow.BinaryTypes.LargeBinary,
		},
		{
			name: "Test OCIBlobLocator",
			args: args{typ: "OCIBlobLocator"},
			want: arrow.BinaryTypes.LargeBinary,
		},
		{
			name: "Test DATE",
			args: args{typ: "DATE"},
			want: arrow.FixedWidthTypes.Date32,
		},
		{
			name: "Test TimeStampDTY",
			args: args{typ: "TimeStampDTY"},
			want: &arrow.TimestampType{Unit: arrow.Microsecond, TimeZone: ""},
		},
		{
			name: "Test TimeStampTZ_DTY",
			args: args{typ: "TimeStampTZ_DTY"},
			want: arrow.FixedWidthTypes.Timestamp_us,
		},
		{
			name: "Test IntervalYM_DTY",
			args: args{typ: "IntervalYM_DTY"},
			want: &arrow.MonthIntervalType{},
		},
		{
			name: "Test IntervalDS_DTY",
			args: args{typ: "IntervalDS_DTY"},
			want: &arrow.DurationType{Unit: arrow.TimeUnit(arrow.Second)},
		},
		// {
		// 	name: "Test JSON",
		// 	want: &extension.JSONType{},
		// },
		{
			name: "Test Unknown Type",
			args: args{typ: "UNKNOWN"},
			want: arrow.BinaryTypes.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := To(tt.args.typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "Test NCHAR OID",
			args: args{oid: 96}, // Assuming 96 corresponds to NCHAR
			want: arrow.BinaryTypes.String,
		},
		{
			name: "Test NUMBER OID",
			args: args{oid: 2}, // Assuming 2 corresponds to NUMBER
			want: arrow.PrimitiveTypes.Float64,
		},
		{
			name: "Test DATE OID",
			args: args{oid: 12}, // Assuming 12 corresponds to DATE
			want: arrow.FixedWidthTypes.Date32,
		},
		{
			name: "Test RAW OID",
			args: args{oid: 23}, // Assuming 23 corresponds to RAW
			want: arrow.BinaryTypes.Binary,
		},
		{
			name: "Test Unknown OID",
			args: args{oid: 999}, // Unknown OID
			want: arrow.BinaryTypes.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToForOID(tt.args.oid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToForOID() = %v, want %v, args.oid = %v", got, tt.want, tt.args.oid)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	tests := []struct {
		name string
		typ  arrow.DataType
		want string
	}{
		{
			name: "Test String",
			typ:  arrow.BinaryTypes.String,
			want: "LONG",
		},
		{
			name: "Test LargeString",
			typ:  arrow.BinaryTypes.LargeString,
			want: "LONG",
		},
		{
			name: "Test Binary",
			typ:  arrow.BinaryTypes.Binary,
			want: "LONG",
		},
		{
			name: "Test LargeBinary",
			typ:  arrow.BinaryTypes.LargeBinary,
			want: "LONG",
		},
		{
			name: "Test Float64",
			typ:  arrow.PrimitiveTypes.Float64,
			want: "LONG",
		},
		{
			name: "Test Float32",
			typ:  arrow.PrimitiveTypes.Float32,
			want: "LONG",
		},
		{
			name: "Test Date32",
			typ:  arrow.FixedWidthTypes.Date32,
			want: "LONG",
		},
		{
			name: "Test Timestamp_us",
			typ:  arrow.FixedWidthTypes.Timestamp_us,
			want: "LONG",
		},
		{
			name: "Test MonthIntervalType",
			typ:  &arrow.MonthIntervalType{},
			want: "LONG",
		},
		{
			name: "Test DurationType",
			typ:  &arrow.DurationType{Unit: arrow.TimeUnit(arrow.Second)},
			want: "LONG",
		},
		// {
		// 	name: "Test JSONType",
		// 	typ:  &extensions.JSONType{},
		// 	want: "LONG",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From(tt.typ); got != tt.want {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}
}
