package client

import (
	"context"
	"errors"

	"github.com/google/wire"
	go_ora "github.com/sijms/go-ora/v2"

	"github.com/openhdc/openhdc"
)

var ProviderSet = wire.NewSet(NewConnector, openhdc.NewDefaultCodec)

type Client struct {
	openhdc.Codec
	opts options

	pool *go_ora.Connection
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.connString == "" {
		return nil, errors.New("connection string is empty")
	}

	f, err := go_ora.ParseConfig(opts.connString)
	if err != nil {
		return nil, err
	}

	p, err := go_ora.NewConnection(opts.connString, f)
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
	return c.pool.Ping(ctx)
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	c.pool.Close()
	return nil
}
