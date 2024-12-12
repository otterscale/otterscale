package connector

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/codec"
)

var ErrNotImplemented = status.Errorf(codes.Unimplemented, "not implemented")

type Connector interface {
	codec.Codec

	Read(ctx context.Context, msg chan<- *pb.Message, opts ReadOptions) error
	Write(ctx context.Context, msg <-chan *pb.Message, opts WriteOptions) error
	Close(ctx context.Context) error
}
