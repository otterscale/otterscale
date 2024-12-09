package client

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/internal/adapter"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

type Client struct {
	opts options
	pool *pgxpool.Pool
}

func NewAdapter(opts ...Option) (adapter.Adapter, error) {
	o := options{
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
		batchTimeout:   defaultBatchTimeout,
		createIndex:    true,
	}
	for _, opt := range opts {
		opt(&o)
	}

	var err error
	o.config, err = pgxpool.ParseConfig(o.connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, o.config)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts: o,
		pool: pool,
	}, nil
}

func (c *Client) TestConnection(ctx context.Context) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
