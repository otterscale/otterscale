package client

import (
	"time"

	"github.com/openhdc/openhdc/api/property/v1"
)

type Option func(*options)

type options struct {
	syncMode       property.SyncMode
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	filePath       string
	tableName      string
	infering       bool
}

func WithSyncMode(syncMode property.SyncMode) Option {
	return func(o *options) {
		o.syncMode = syncMode
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
