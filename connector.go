package openhdc

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Connector interface {
	Codec
	Writer

	Read(ctx context.Context, msgs chan<- *pb.Message, opts ReadOptions) error
	Write(ctx context.Context, msgs <-chan *pb.Message) error
	Close(ctx context.Context) error
}
