package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"

	"github.com/otterscale/otterscale/internal/core"
)

func TestDomainErrorToConnectError_ConcreteTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode connect.Code
	}{
		{
			name:     "ErrInvalidInput",
			err:      &core.ErrInvalidInput{Field: "name", Message: "required"},
			wantCode: connect.CodeInvalidArgument,
		},
		{
			name:     "ErrSessionNotFound",
			err:      &core.ErrSessionNotFound{Resource: "exec", ID: "abc"},
			wantCode: connect.CodeNotFound,
		},
		{
			name:     "ErrClusterNotFound",
			err:      &core.ErrClusterNotFound{Cluster: "test"},
			wantCode: connect.CodeNotFound,
		},
		{
			name:     "ErrNotReady",
			err:      &core.ErrNotReady{Subsystem: "chisel"},
			wantCode: connect.CodeUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domainErrorToConnectError(tt.err)
			var connectErr *connect.Error
			if !errors.As(got, &connectErr) {
				t.Fatalf("expected *connect.Error, got %T", got)
			}
			if connectErr.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, connectErr.Code())
			}
		})
	}
}

func TestDomainErrorToConnectError_DomainErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     core.ErrorCode
		wantCode connect.Code
	}{
		{"Internal", core.ErrorCodeInternal, connect.CodeInternal},
		{"InvalidArgument", core.ErrorCodeInvalidArgument, connect.CodeInvalidArgument},
		{"NotFound", core.ErrorCodeNotFound, connect.CodeNotFound},
		{"AlreadyExists", core.ErrorCodeAlreadyExists, connect.CodeAlreadyExists},
		{"Unauthenticated", core.ErrorCodeUnauthenticated, connect.CodeUnauthenticated},
		{"PermissionDenied", core.ErrorCodePermissionDenied, connect.CodePermissionDenied},
		{"FailedPrecondition", core.ErrorCodeFailedPrecondition, connect.CodeFailedPrecondition},
		{"DeadlineExceeded", core.ErrorCodeDeadlineExceeded, connect.CodeDeadlineExceeded},
		{"ResourceExhausted", core.ErrorCodeResourceExhausted, connect.CodeResourceExhausted},
		{"Unimplemented", core.ErrorCodeUnimplemented, connect.CodeUnimplemented},
		{"Unavailable", core.ErrorCodeUnavailable, connect.CodeUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &core.DomainError{Code: tt.code, Message: "test"}
			got := domainErrorToConnectError(err)
			var connectErr *connect.Error
			if !errors.As(got, &connectErr) {
				t.Fatalf("expected *connect.Error, got %T", got)
			}
			if connectErr.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, connectErr.Code())
			}
		})
	}
}

func TestDomainErrorToConnectError_UnknownError(t *testing.T) {
	got := domainErrorToConnectError(errors.New("random error"))
	var connectErr *connect.Error
	if !errors.As(got, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T", got)
	}
	if connectErr.Code() != connect.CodeInternal {
		t.Errorf("expected CodeInternal for unknown error, got %v", connectErr.Code())
	}
}

func TestDomainCodeToConnectCode_Completeness(t *testing.T) {
	// Verify the map has entries for all defined error codes.
	if len(domainCodeToConnectCode) < 11 {
		t.Errorf("expected at least 11 domain code mappings, got %d", len(domainCodeToConnectCode))
	}
}

// TestDomainErrorToConnectError_ContextErrors covers the codes a client
// that hangs up should produce. Before these cases existed, a bare
// ctx.Err() traveling up from a use-case was reported as an internal
// fault, so every abandoned request looked like a server bug.
func TestDomainErrorToConnectError_ContextErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode connect.Code
	}{
		{
			name:     "canceled",
			err:      context.Canceled,
			wantCode: connect.CodeCanceled,
		},
		{
			name:     "deadline exceeded",
			err:      context.DeadlineExceeded,
			wantCode: connect.CodeDeadlineExceeded,
		},
		{
			name:     "wrapped cancellation",
			err:      fmt.Errorf("read pod logs: %w", context.Canceled),
			wantCode: connect.CodeCanceled,
		},
		{
			// An adapter that classified the cancellation itself keeps
			// the code it chose: the context check runs after the domain
			// types precisely so explicit intent wins.
			name: "domain code wins over context cause",
			err: &core.DomainError{
				Code:    core.ErrorCodeDeadlineExceeded,
				Message: "gave up fetching chart",
				Cause:   context.Canceled,
			},
			wantCode: connect.CodeDeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := connect.CodeOf(domainErrorToConnectError(tt.err)); got != tt.wantCode {
				t.Errorf("code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

// TestDomainErrorToConnectError_SchemaNotFound pins the mapping for a
// kind the cluster does not define. Asking for an unknown kind is a bad
// request; it used to surface as an internal error because the adapter
// returned the upstream sentinel without a domain code.
func TestDomainErrorToConnectError_SchemaNotFound(t *testing.T) {
	err := &core.DomainError{
		Code:    core.ErrorCodeNotFound,
		Message: `cannot resolve group version kind "apps/v1/Nonexistent"`,
		Cause:   resolver.ErrSchemaNotFound,
	}

	if got := connect.CodeOf(domainErrorToConnectError(err)); got != connect.CodeNotFound {
		t.Errorf("code = %v, want %v", got, connect.CodeNotFound)
	}
	if !errors.Is(err, resolver.ErrSchemaNotFound) {
		t.Error("the upstream sentinel must stay in the chain for errors.Is")
	}
}
