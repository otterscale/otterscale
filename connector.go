package openhdc

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Connector interface {
	Codec

	Read(ctx context.Context, msg chan<- *pb.Message, opts ReadOptions) error
	Write(ctx context.Context, msg <-chan *pb.Message, opts WriteOptions) error
	Close(ctx context.Context) error
}
