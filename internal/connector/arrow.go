package connector

import (
	"bytes"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

func FromRecord(rec arrow.Record) ([]byte, error) {
	var buf bytes.Buffer
	w := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()), ipc.WithAllocator(memory.DefaultAllocator))
	defer w.Close()
	if err := w.Write(rec); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ToRecord(b []byte) (arrow.Record, error) {
	r, err := ipc.NewReader(bytes.NewReader(b), ipc.WithAllocator(memory.DefaultAllocator))
	if err != nil {
		return nil, err
	}
	defer r.Release()
	return r.Record(), nil
}
