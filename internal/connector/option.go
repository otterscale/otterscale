package connector

import (
	"google.golang.org/grpc"
)

type Option func(*options)

type options struct {
	kind Kind
}

func WithKind(kind Kind) Option {
	return func(o *options) {
		o.kind = kind
	}
}

type ServerOption func(*serverOptions)

type serverOptions struct {
	network  string
	address  string
	grpcOpts []grpc.ServerOption
}

func WithNetwork(network string) ServerOption {
	return func(o *serverOptions) {
		o.network = network
	}
}

func WithAddress(address string) ServerOption {
	return func(o *serverOptions) {
		o.address = address
	}
}

func WithGrpcServerOptions(grpcOpts ...grpc.ServerOption) ServerOption {
	return func(o *serverOptions) {
		o.grpcOpts = append(o.grpcOpts, grpcOpts...)
	}
}

type ReadOption func(*readOptions)

type readOptions struct {
	tables     []string
	skipTables []string
}

func WithTables(tables []string) ReadOption {
	return func(o *readOptions) {
		o.tables = tables
	}
}

func WithSkipTables(skipTables []string) ReadOption {
	return func(o *readOptions) {
		o.skipTables = skipTables
	}
}

type WriteOption func(*writeOptions)

type writeOptions struct {
	kind  WriteKind
	table string
}

func WithWriteKind(kind WriteKind) WriteOption {
	return func(o *writeOptions) {
		o.kind = kind
	}
}

func WithTable(table string) WriteOption {
	return func(o *writeOptions) {
		o.table = table
	}
}
