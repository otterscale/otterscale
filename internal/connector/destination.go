package connector

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Destination interface {
	Write(ctx context.Context, record chan<- arrow.Record) error
	Close(ctx context.Context) error
}

type UnimplementedDestination struct{}

func (UnimplementedDestination) Write(ctx context.Context, record chan<- arrow.Record) error {
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}

func (UnimplementedDestination) Close(ctx context.Context) error {
	return status.Errorf(codes.Unimplemented, "method Close not implemented")
}
