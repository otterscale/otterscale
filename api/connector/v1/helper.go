package pb

import (
	"bytes"
	"errors"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openhdc/openhdc/api/property/v1"
)

func NewMessage(k property.MessageKind, r arrow.Record, sourceName string, syncedAt time.Time) (*Message, error) {
	var buf bytes.Buffer
	w := ipc.NewWriter(&buf, ipc.WithSchema(r.Schema()), ipc.WithAllocator(memory.DefaultAllocator))
	defer w.Close()
	if err := w.Write(r); err != nil {
		return nil, err
	}
	t := timestamppb.New(syncedAt)
	return Message_builder{
		Kind:       &k,
		Record:     buf.Bytes(),
		SourceName: &sourceName,
		SyncedAt:   t,
	}.Build(), nil
}

func ToRecord(b []byte) (arrow.Record, error) {
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
