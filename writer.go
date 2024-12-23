package openhdc

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Writer interface {
	Write(ctx context.Context, msgs <-chan *pb.Message) error
}

func NewEmptyWriter() Writer {
	return &EmptyWriter{}
}

type EmptyWriter struct{}

func (EmptyWriter) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	return ErrNotImplemented
}
