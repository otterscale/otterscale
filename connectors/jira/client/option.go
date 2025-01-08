package client

import (
	"time"
)

type Option func(*options)

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

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithServer(server string) Option {
	return func(o *options) {
		o.server = server
	}
}

func WithUsername(username string) Option {
	return func(o *options) {
		o.username = username
	}
}

func WithPassword(password string) Option {
	return func(o *options) {
		o.password = password
	}
}

func WithToken(token string) Option {
	return func(o *options) {
		o.token = token
	}
}

func WithProjects(projects []string) Option {
	return func(o *options) {
		o.projects = projects
	}
}

func WithStartDate(startDate string) Option {
	return func(o *options) {
		o.startDate = startDate
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
