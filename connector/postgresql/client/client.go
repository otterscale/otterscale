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

	config *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewConnector(c openhdc.Codec, opts ...Option) (openhdc.Connector, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.connString == "" {
		return nil, errors.New("connection string is empty")
	}

	var err error
	f, err := pgxpool.ParseConfig(o.connString)
	if err != nil {
		return nil, err
	}

	p, err := pgxpool.NewWithConfig(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return &Client{
		Codec:  c,
		opts:   o,
		config: f,
		pool:   p,
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
