package adapter

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
)

type Adapter interface {
	GetTables(ctx context.Context) ([]arrow.Table, error)
	Migrate(ctx context.Context)

	List(ctx context.Context, record chan<- arrow.Record) error
	Create(ctx context.Context, record chan<- arrow.Record) error
	Delete(ctx context.Context)
}
