package client

import (
	"context"
	"errors"
	"os"

	"github.com/openhdc/openhdc"
)

type Client struct {
	openhdc.Codec
	opts options

	file *os.File
}

func NewConnector(c openhdc.Codec, opt ...Option) (openhdc.Connector, error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	if opts.filePath == "" {
		return nil, errors.New("file path is empty")
	}

	f, err := os.Open(opts.filePath)
	if err != nil {
		return nil, err
	}

	return &Client{
		Codec: c,
		opts:  opts,
		file:  f,
	}, nil
}

func (c *Client) Name() string {
	return c.opts.name
}

func (c *Client) Close(ctx context.Context) error {
	return c.file.Close()
}
