package client

import (
	"context"

	"github.com/openhdc/openhdc"
	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Writer struct{}

func NewWriter() openhdc.Writer {
	return &Writer{}
}

func (Writer) Write(ctx context.Context, msgs <-chan *pb.Message) error {
	return openhdc.ErrNotImplemented
}
