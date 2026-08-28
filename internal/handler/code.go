package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/core"
)

// domainCodeToConnectCode maps domain-level error codes to their
// ConnectRPC equivalents.
var domainCodeToConnectCode = map[core.ErrorCode]connect.Code{
	core.ErrorCodeInternal:           connect.CodeInternal,
	core.ErrorCodeInvalidArgument:    connect.CodeInvalidArgument,
	core.ErrorCodeNotFound:           connect.CodeNotFound,
	core.ErrorCodeAlreadyExists:      connect.CodeAlreadyExists,
	core.ErrorCodeUnauthenticated:    connect.CodeUnauthenticated,
	core.ErrorCodePermissionDenied:   connect.CodePermissionDenied,
	core.ErrorCodeFailedPrecondition: connect.CodeFailedPrecondition,
	core.ErrorCodeDeadlineExceeded:   connect.CodeDeadlineExceeded,
	core.ErrorCodeResourceExhausted:  connect.CodeResourceExhausted,
	core.ErrorCodeUnimplemented:      connect.CodeUnimplemented,
	core.ErrorCodeUnavailable:        connect.CodeUnavailable,
	core.ErrorCodeCanceled:           connect.CodeCanceled,
}

// domainErrorToConnectError converts a domain error into a ConnectRPC
// error with a semantically equivalent code. Domain-specific error
// types (ErrInvalidInput, ErrClusterNotFound, etc.) are checked first,
// then DomainError codes are mapped. Unrecognized errors fall back to
// connect.CodeInternal.
func domainErrorToConnectError(err error) error {
	// Concrete domain error types.
	var invalidInput *core.ErrInvalidInput
	if errors.As(err, &invalidInput) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	var sessionNotFound *core.ErrSessionNotFound
	if errors.As(err, &sessionNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	var clusterNotFound *core.ErrClusterNotFound
	if errors.As(err, &clusterNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	var notReady *core.ErrNotReady
	if errors.As(err, &notReady) {
		return connect.NewError(connect.CodeUnavailable, err)
	}

	// Generic domain error with error code.
	var domainErr *core.DomainError
	if errors.As(err, &domainErr) {
		code, ok := domainCodeToConnectCode[domainErr.Code]
		if !ok {
			code = connect.CodeInternal
		}
		return connect.NewError(code, err)
	}

	// Context errors are checked after the domain types so that an
	// adapter which deliberately classified a cancellation keeps the
	// code it chose. What lands here is a bare ctx.Err() traveling up
	// from a use-case or from client-go, which would otherwise be
	// reported as an internal fault and logged as a server error —
	// every client that hangs up mid-request would look like a bug.
	if errors.Is(err, context.Canceled) {
		return connect.NewError(connect.CodeCanceled, err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return connect.NewError(connect.CodeDeadlineExceeded, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}
