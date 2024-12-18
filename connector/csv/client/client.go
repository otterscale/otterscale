package client

import (
	"context"
	"log/slog"
	"os"

	"github.com/openhdc/openhdc"
)

type Client struct {
	openhdc.Codec
	opts options

	file *os.File
}

func NewConnector(opts ...Option) (openhdc.Connector, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	f, err := os.Open(o.path)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &Client{
		Codec: openhdc.DefaultCodec{},
		opts:  o,
		file:  f,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	if err := c.file.Close(); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
