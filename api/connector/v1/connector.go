package pb

import (
	"bytes"
	"errors"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openhdc/openhdc/api/property/v1"
)

var ErrNotSupported = status.Errorf(codes.InvalidArgument, "not supported")

func NewMessage(kind property.MessageKind, rec arrow.Record, sourceName string, syncedAt time.Time) (*Message, error) {
	b, err := FromArrowRecord(rec)
	if err != nil {
		return nil, err
	}
	m := &Message{}
	m.SetKind(kind)
	m.SetRecord(b)
	m.SetSourceName(sourceName)
	m.SetSyncedAt(timestamppb.New(syncedAt))
	return m, nil
}

func ToArrowRecord(b []byte) (arrow.Record, error) {
	r, err := ipc.NewReader(bytes.NewReader(b), ipc.WithAllocator(memory.DefaultAllocator))
	if err != nil {
		return nil, err
	}
	defer r.Release()
	for r.Next() {
		rec := r.Record()
		rec.Retain()
		return rec, nil
	}
	return nil, errors.New("record is empty")
}

func FromArrowRecord(rec arrow.Record) ([]byte, error) {
	var buf bytes.Buffer
	w := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()), ipc.WithAllocator(memory.DefaultAllocator))
	defer w.Close()
	if err := w.Write(rec); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
