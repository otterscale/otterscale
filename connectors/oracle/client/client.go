package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/wire"

	"github.com/openhdc/openhdc"
)

var ProviderSet = wire.NewSet(NewConnector, openhdc.NewDefaultCodec)

type Client struct {
	openhdc.Codec
	opts options

	pool *sql.DB
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.connString == "" {
		return nil, errors.New("connection string is empty")
	}

	u, err := url.Parse(opts.connString)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.Scheme)

	p, err := sql.Open(u.Scheme, opts.connString)
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
	return c.pool.PingContext(ctx)
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	c.pool.Close()
	return nil
}
