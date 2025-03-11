package client

import (
	"time"
)

type options struct {
	name           string
	server         string
	username       string
	password       string
	token          string
	projects       []string
	startDate      string
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

func WithName(t string) Option {
	return newFuncOption(func(o *options) {
		o.name = t
	})
}

func WithServer(t string) Option {
	return newFuncOption(func(o *options) {
		o.server = t
	})
}

func WithUsername(t string) Option {
	return newFuncOption(func(o *options) {
		o.username = t
	})
}

func WithPassword(t string) Option {
	return newFuncOption(func(o *options) {
		o.password = t
	})
}

func WithToken(t string) Option {
	return newFuncOption(func(o *options) {
		o.token = t
	})
}

func WithProjects(tl []string) Option {
	return newFuncOption(func(o *options) {
		o.projects = tl
	})
}

func WithStartDate(t string) Option {
	return newFuncOption(func(o *options) {
		o.startDate = t
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
