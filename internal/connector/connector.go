package connector

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openhdc/openhdc/internal/codec"
)

var ErrNotImplemented = status.Errorf(codes.Unimplemented, "not implemented")

type Connector interface {
	codec.Codec

	Read(ctx context.Context, rec chan<- arrow.Record, opts ...ReadOption) error
	Write(ctx context.Context, rec chan<- arrow.Record, opts ...WriteOption) error
	Close(ctx context.Context) error
}
