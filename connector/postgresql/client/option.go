package client

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*options)

type options struct {
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	createIndex    bool

	connString string
	config     *pgxpool.Config
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

func WithCreateIndex(createIndex bool) Option {
	return func(o *options) {
		o.createIndex = createIndex
	}
}

func WithConnConfig(connString string) Option {
	return func(o *options) {
		o.connString = connString
	}
}
