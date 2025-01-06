package openhdc

import (
	"slices"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

const (
	builtinFieldName     = "_openhdc_name"
	builtinFieldSyncedAt = "_openhdc_synced_at"
)

func BuiltinFieldName() string {
	return builtinFieldName
}

func BuiltinFieldSyncedAt() string {
	return builtinFieldSyncedAt
}

func isBuiltinField(a *arrow.Field) bool {
	return slices.Contains([]string{builtinFieldName, builtinFieldSyncedAt}, a.Name)
}

func appendBuiltinFields(rec arrow.Record) *arrow.Schema {
	fs := []arrow.Field{
		{
			Name:     builtinFieldName,
			Type:     arrow.BinaryTypes.String,
			Nullable: false,
		},
		{
			Name:     builtinFieldSyncedAt,
			Type:     arrow.FixedWidthTypes.Timestamp_us,
			Nullable: false,
		},
	}
	md := rec.Schema().Metadata()
	return arrow.NewSchema(append(fs, rec.Schema().Fields()...), &md)
}

func appendBuiltinValues(rec arrow.Record, name string, syncedAt time.Time) []arrow.Array {
	sb := array.NewStringBuilder(memory.DefaultAllocator)
	tb := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_us.(*arrow.TimestampType))
	for range rec.NumRows() {
		sb.Append(name)
		tb.AppendTime(syncedAt)
	}
	arrs := make([]arrow.Array, rec.NumCols())
	copy(arrs, rec.Columns())
	return append([]arrow.Array{sb.NewArray(), tb.NewArray()}, arrs...)
}

func AppendBuiltinFieldsToRecord(msg *pb.Message) (arrow.Record, error) {
	rec, err := pb.ToRecord(msg.GetRecord())
	if err != nil {
		return nil, err
	}
	sch := appendBuiltinFields(rec)
	arrs := appendBuiltinValues(rec, msg.GetSourceName(), msg.GetSyncedAt().AsTime())
	return array.NewRecord(sch, arrs, rec.NumRows()), nil
}
