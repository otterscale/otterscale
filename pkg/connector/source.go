package connector

import (
	"context"
)

type Source interface {
	Read(ctx context.Context) error
	Close(ctx context.Context) error
}
