package client

import (
	"time"
)

type Option func(*options)

type options struct {
	name           string
	filePath       string
	tableName      string
	infering       bool
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	createIndex    bool
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithFilePath(filePath string) Option {
	return func(o *options) {
		o.filePath = filePath
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
