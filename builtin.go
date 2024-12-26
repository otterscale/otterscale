package openhdc

import (
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

const (
	BuiltinFieldName     = "_openhdc_name"
	BuiltinFieldSyncedAt = "_openhdc_synced_at"
)

func newSchema(rec arrow.Record) *arrow.Schema {
	fs := append([]arrow.Field{
		{Name: BuiltinFieldName, Type: arrow.BinaryTypes.String, Nullable: false},
		{Name: BuiltinFieldSyncedAt, Type: arrow.FixedWidthTypes.Timestamp_us, Nullable: false},
	}, rec.Schema().Fields()...)
	md := rec.Schema().Metadata()
	return arrow.NewSchema(fs, &md)
}

func newColumns(rec arrow.Record, name string, syncedAt time.Time) []arrow.Array {
	sb := array.NewStringBuilder(memory.DefaultAllocator)
	tb := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_us.(*arrow.TimestampType))
	for range rec.NumRows() {
		sb.Append(name)
		tb.AppendTime(syncedAt)
	}
	arrs := make([]arrow.Array, rec.NumCols())
	copy(arrs, rec.Columns())
	return append([]arrow.Array{
		sb.NewArray(),
		tb.NewArray(),
	}, arrs...)
}

func AppendBuiltinFieldsToRecord(msg *pb.Message) (arrow.Record, error) {
	rec, err := pb.ToArrowRecord(msg.GetRecord())
	if err != nil {
		return nil, err
	}
	sch := newSchema(rec)
	cols := newColumns(rec, msg.GetSourceName(), msg.GetSyncedAt().AsTime())
	return array.NewRecord(sch, cols, rec.NumRows()), nil
}
