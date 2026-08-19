package handler

import (
	"errors"
	"testing"

	"connectrpc.com/connect"

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
