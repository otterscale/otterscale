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

func NewConnector(opts ...Option) (openhdc.Connector, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.filePath == "" {
		return nil, errors.New("file path is empty")
	}

	f, err := os.Open(o.filePath)
	if err != nil {
		return nil, err
	}

	return &Client{
		Codec: openhdc.DefaultCodec{},
		opts:  o,
		file:  f,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.file.Close()
}
