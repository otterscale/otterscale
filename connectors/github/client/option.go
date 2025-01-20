package client

import (
	"time"
)

type options struct {
	name           string
	owner          string
	repo           string
	opts           string
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

func WithOwner(ow string) Option {
	return newFuncOption(func(o *options) {
		o.owner = ow
	})
}

func WithRepo(r string) Option {
	return newFuncOption(func(o *options) {
		o.repo = r
	})
}

func WithOpts(op string) Option {
	return newFuncOption(func(o *options) {
		o.opts = op
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
