package connector

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
)

type Connector interface {
	Read(ctx context.Context, record chan<- arrow.Record, opts ...ReadOption) error
	Write(ctx context.Context, record chan<- arrow.Record) error
	Close(ctx context.Context) error
}
