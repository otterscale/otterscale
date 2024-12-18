package client

import (
	"time"
)

type Option func(*options)

type options struct {
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration

	path      string
	tableName string

	infering bool
}

func WithBatchSize(batchSize int64) Option {
	return func(o *options) {
		o.batchSize = batchSize
	}
}

func WithBatchSizeBytes(batchSizeBytes int64) Option {
	return func(o *options) {
		o.batchSizeBytes = batchSizeBytes
	}
}

func WithBatchTimeout(batchTimeout time.Duration) Option {
	return func(o *options) {
		o.batchTimeout = batchTimeout
	}
}

func WithPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}

func WithTableName(tableName string) Option {
	return func(o *options) {
		o.tableName = tableName
	}
}

func WithInfering(infering bool) Option {
	return func(o *options) {
		o.infering = infering
	}
}
