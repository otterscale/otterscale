package client

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc/internal/codec"
	"github.com/openhdc/openhdc/internal/connector"
)

const (
	defaultBatchSize      = 10000
	defaultBatchSizeBytes = 100000000
	defaultBatchTimeout   = 60 * time.Second
)

type Client struct {
	codec.Codec
	opts options

	config *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewConnector(codec codec.Codec, opts ...Option) (connector.Connector, error) {
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
	config, err := pgxpool.ParseConfig(o.connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Codec:  codec,
		opts:   o,
		config: config,
		pool:   pool,
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

func (c *Client) Close(ctx context.Context) error {
	c.pool.Close()
	return nil
}
