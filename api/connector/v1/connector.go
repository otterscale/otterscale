package pb

import (
	"bytes"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotSupported = status.Errorf(codes.InvalidArgument, "not supported")

func NewMessage(kind Kind, rec arrow.Record) (*Message, error) {
	b, err := FromArrowRecord(rec)
	if err != nil {
		return nil, err
	}
	return &Message{
		Kind:   kind,
		Record: b,
	}, nil
}

func ToArrowRecord(b []byte) (arrow.Record, error) {
	r, err := ipc.NewReader(bytes.NewReader(b), ipc.WithAllocator(memory.DefaultAllocator))
	if err != nil {
		return nil, err
	}
	defer r.Release()
	return r.Read()
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
