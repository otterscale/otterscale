package connector

import (
	"context"
)

type Destination interface {
	Write(ctx context.Context) error
	Close(ctx context.Context) error
}
