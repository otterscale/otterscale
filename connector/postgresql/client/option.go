package client

import (
	"time"

	"github.com/openhdc/openhdc/api/property/v1"
)

type Option func(*options)

type options struct {
	name           string
	syncMode       property.SyncMode
	cursor         string
	connString     string
	namespace      string
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

func WithSyncMode(syncMode property.SyncMode) Option {
	return func(o *options) {
		o.syncMode = syncMode
	}
}

func WithCursor(cursor string) Option {
	return func(o *options) {
		o.cursor = cursor
	}
}

func WithConnString(connString string) Option {
	return func(o *options) {
		o.connString = connString
	}
}

func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
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
