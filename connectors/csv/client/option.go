package client

import (
	"time"
)

type options struct {
	name           string
	filePath       string
	tableName      string
	inferring      bool
	batchSize      int64
	batchSizeBytes int64
	batchTimeout   time.Duration
	createIndex    bool
}

var defaultOptions = options{}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

var _ Option = (*funcOption)(nil)

func (fro *funcOption) apply(ro *options) {
	fro.f(ro)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithName(n string) Option {
	return newFuncOption(func(o *options) {
		o.name = n
	})
}

func WithFilePath(f string) Option {
	return newFuncOption(func(o *options) {
		o.filePath = f
	})
}

func WithTableName(t string) Option {
	return newFuncOption(func(o *options) {
		o.tableName = t
	})
}

func WithInferring(b bool) Option {
	return newFuncOption(func(o *options) {
		o.inferring = b
	})
}

func WithBatchSize(s int64) Option {
	return newFuncOption(func(o *options) {
		o.batchSize = s
	})
}

func WithBatchSizeBytes(s int64) Option {
	return newFuncOption(func(o *options) {
		o.batchSizeBytes = s
	})
}

func WithBatchTimeout(t time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.batchTimeout = t
	})
}

func WithCreateIndex(b bool) Option {
	return newFuncOption(func(o *options) {
		o.createIndex = b
	})
}
