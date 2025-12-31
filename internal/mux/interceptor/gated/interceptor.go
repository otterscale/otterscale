package gated

import (
	"context"

	"connectrpc.com/connect"
)

type Interceptor struct {
	interceptor connect.Interceptor
	shouldRun   func(ctx context.Context) bool
}

func NewInterceptor(interceptor connect.Interceptor, shouldRun func(ctx context.Context) bool) *Interceptor {
	return &Interceptor{
		interceptor: interceptor,
		shouldRun:   shouldRun,
	}
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if i.shouldRun(ctx) {
			return i.interceptor.WrapUnary(next)(ctx, req)
		}
		return next(ctx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		if i.shouldRun(ctx) {
			return i.interceptor.WrapStreamingClient(next)(ctx, spec)
		}
		return next(ctx, spec)
	}
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if i.shouldRun(ctx) {
			return i.interceptor.WrapStreamingHandler(next)(ctx, conn)
		}
		return next(ctx, conn)
	}
}
