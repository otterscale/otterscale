package connector

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Source interface {
	Read(ctx context.Context, record chan<- arrow.Record) error
	Close(ctx context.Context) error
}

type UnimplementedSource struct{}

func (UnimplementedSource) Read(ctx context.Context, record chan<- arrow.Record) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}

func (UnimplementedSource) Close(ctx context.Context) error {
	return status.Errorf(codes.Unimplemented, "method Close not implemented")
}
