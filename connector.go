package openhdc

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
)

type Connector interface {
	Codec

	Name() string
	Read(ctx context.Context, msgs chan<- *pb.Message, rdr *Reader) error
	Write(ctx context.Context, msgs <-chan *pb.Message) error
	Close(ctx context.Context) error
}
