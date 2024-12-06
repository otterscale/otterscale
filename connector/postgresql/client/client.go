package client

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/pkg/adapter"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

type Client struct {
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	createIndex    bool

	connString string
	config     *pgxpool.Config
	pool       *pgxpool.Pool
}

func NewAdapter(opts ...Option) (adapter.Adapter, error) {
	c := &Client{
		batchSize:      defaultBatchSize,
		batchSizeBytes: defaultBatchSizeBytes,
		batchTimeout:   defaultBatchTimeout,
		createIndex:    true,
	}
	for _, opt := range opts {
		opt(c)
	}

	var err error
	c.config, err = pgxpool.ParseConfig(c.connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.pool, err = pgxpool.NewWithConfig(ctx, c.config)
	if err != nil {
		return nil, err
	}

	return c, nil
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
