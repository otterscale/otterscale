package client

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/openhdc/openhdc"
)

type Client struct {
	openhdc.Codec
	opts options

	pool *pgxpool.Pool
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.connString == "" {
		return nil, errors.New("connection string is empty")
	}

	f, err := pgxpool.ParseConfig(opts.connString)
	if err != nil {
		return nil, err
	}

	p, err := pgxpool.NewWithConfig(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return &Client{
		Codec: c,
		opts:  opts,
		pool:  p,
	}, nil
}

func (c *Client) TestConnection(ctx context.Context) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.Ping(ctx)
}

func (c *Client) Close(ctx context.Context) error {
	c.pool.Close()
	return nil
}
