package client

import "time"

type Option func(*Client)

func WithBatchSize(batchSize int64) Option {
	return func(c *Client) {
		c.batchSize = batchSize
	}
}

func WithBatchSizeBytes(batchSizeBytes int64) Option {
	return func(c *Client) {
		c.batchSizeBytes = batchSizeBytes
	}
}

func WithBatchTimeout(batchTimeout time.Duration) Option {
	return func(c *Client) {
		c.batchTimeout = batchTimeout
	}
}

func WithCreateIndex(createIndex bool) Option {
	return func(c *Client) {
		c.createIndex = createIndex
	}
}

func WithConnConfig(connString string) Option {
	return func(c *Client) {
		c.connString = connString
	}
}
